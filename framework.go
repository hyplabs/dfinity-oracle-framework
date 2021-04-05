package framework

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyplabs/dfinity-oracles/models"
	"github.com/hyplabs/dfinity-oracles/summary"
	"github.com/hyplabs/dfinity-oracles/utils"
	"github.com/sirupsen/logrus"
)

// Oracle is an instance of an oracle
type Oracle struct {
	config     *models.Config
	dfxService *DFXService
	engine     *models.Engine
	log        *logrus.Logger
}

// NewOracle creates a new oracle instance
func NewOracle(config *models.Config, engine *models.Engine) *Oracle {
	log := logrus.New()
	log.Formatter = &logrus.JSONFormatter{}

	dfxService := NewDFXService(config, log)

	return &Oracle{
		config:     config,
		dfxService: dfxService,
		engine:     engine,
		log:        log,
	}
}

// Bootstrap bootstraps the canister installation
func (o *Oracle) Bootstrap() error {
	if err := o.dfxService.createNewDfxProject(); err != nil {
		panic(err)
	}
	if err := o.dfxService.updateCanisterCode(); err != nil {
		panic(err)
	}
	if err := o.dfxService.stopDfxNetwork(); err != nil {
		panic(err)
	}
	if err := o.dfxService.startDfxNetwork(); err != nil {
		panic(err)
	}
	if err := o.dfxService.createWriterIdentityIfNeeded(); err != nil {
		panic(err)
	}

	canisterExists, err := o.dfxService.doesCanisterExist()
	if err != nil {
		panic(err)
	}
	if !canisterExists {
		if err := o.dfxService.createCanister(); err != nil {
			return err
		}
	}
	if err := o.dfxService.buildCanister(); err != nil {
		panic(err)
	}
	if err := o.dfxService.installCanister(canisterExists); err != nil {
		panic(err)
	}
	canisterRunning, err := o.dfxService.isCanisterRunning()
	if err != nil {
		panic(err)
	}
	if !canisterRunning {
		if err := o.dfxService.startCanister(); err != nil {
			panic(err)
		}
	}

	isOwner, err := o.dfxService.checkIsOwner()
	if err != nil {
		panic(err)
	}
	if !isOwner {
		if err := o.dfxService.assignOwnerRole(); err != nil {
			panic(err)
		}
	}
	writerPrincipal, err := o.dfxService.getWriterIDPrincipal()
	if err != nil {
		panic(err)
	}
	if err := o.dfxService.assignWriterRole(writerPrincipal); err != nil {
		panic(err)
	}
	return nil
}

// Run starts the Oracle service
func (o *Oracle) Run() {
	o.log.Infof("Starting %s oracle service...", o.config.CanisterName)

	o.updateOracle()
	for range time.Tick(o.config.UpdateInterval) {
		o.updateOracle()
	}
}

func (o *Oracle) updateOracle() {
	for _, meta := range o.engine.Metadata {
		o.updateMeta(meta)
	}
	o.log.Infof("Oracle update completed")
}

func (o *Oracle) updateMeta(meta models.MappingMetadata) error {
	type apiInfo struct {
		Endpoint models.Endpoint
		Value    map[string]float64
		Err      error
	}
	dataset := make([]map[string]float64, 0)
	ch := make(chan apiInfo)
	for _, endpoint := range meta.Endpoints {
		go func(endpoint models.Endpoint, ch chan<- apiInfo) {
			val, err := utils.GetAPIInfo(endpoint)
			ch <- apiInfo{Endpoint: endpoint, Value: val, Err: err}
		}(endpoint, ch)
	}
	for range meta.Endpoints {
		r := <-ch
		if r.Err != nil {
			o.log.WithError(r.Err).Errorf("Could not retrieve information from API %s", r.Endpoint.Endpoint)
			return r.Err
		}
		dataset = append(dataset, r.Value)
		valStr, err := json.Marshal(r.Value)
		if err != nil {
			o.log.WithError(err).Errorf("Retrieved non-JSON-serializable value %v from %s for %s", r.Value, r.Endpoint.Endpoint, meta.Key)
			return err
		}
		o.log.Infof("Retrieved value %v from %s for %s", string(valStr), r.Endpoint.Endpoint, meta.Key)
	}
	if len(dataset) == 0 {
		o.log.Errorf("No values from any API endpoints, skipping update for %s", meta.Key)
		return fmt.Errorf("No values from any API endpoints, skipping update for %s", meta.Key)
	}

	var summarizedVal map[string]float64
	if meta.SummaryFunc != nil {
		summarizedVal = meta.SummaryFunc(dataset)
	} else {
		summarizedVal = summary.MeanWithoutOutliers(dataset)
	}
	o.dfxService.updateValueInCanister(meta.Key, summarizedVal)
	return nil
}
