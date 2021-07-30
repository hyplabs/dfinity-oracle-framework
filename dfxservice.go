package framework

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hyplabs/dfinity-oracle-framework/models"
	"github.com/hyplabs/dfinity-oracle-framework/utils"
	"github.com/sirupsen/logrus"
)

// DFXService contains various fields to be used by the DFX interface
type DFXService struct {
	config *models.Config
	log    *logrus.Logger
}

// NewDFXService creates an instance of DFX Service
func NewDFXService(config *models.Config, log *logrus.Logger) *DFXService {
	return &DFXService{
		config,
		log,
	}
}

func (s *DFXService) createNewDfxProject() error {
	s.log.Infof("Creating new project %s...", s.config.CanisterName)
	output, exitCode, err := dfxCall(".", []string{"new", s.config.CanisterName}, true)
	if err != nil {
		s.log.WithError(err).Errorln("Could not create new project:", output)
		return err
	}
	if exitCode != 0 && !strings.HasPrefix(output, "Cannot create a new project because the directory already exists.") {
		s.log.WithError(err).Errorln("Could not create new project:", output)
		return fmt.Errorf("Could not create new project: %v", output)
	}
	return nil
}

func (s *DFXService) updateCanisterCode() error {
	s.log.Infof("Updating canister code...")
	fileName := filepath.Join(s.config.CanisterName, "src", s.config.CanisterName, "main.mo")
	if err := ioutil.WriteFile(fileName, []byte(CodeTemplate), 0644); err != nil {
		s.log.WithError(err).Errorln("Could not write to main.mo file")
		return err
	}
	return nil
}

func (s *DFXService) stopDfxNetwork() error {
	s.log.Infof("Stopping existing DFX instances...")
	_, _, err := dfxCall(s.config.CanisterName, []string{"stop"}, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not stop DFX identity")
		return err
	}
	s.log.Infoln("Sleeping 5 seconds to allow local network to finish shutting down...")
	time.Sleep(5 * time.Second)
	return nil
}

func (s *DFXService) startDfxNetwork() error {
	s.log.Infof("Starting DFX in the background...")
	dfxExecutable, err := exec.LookPath("dfx")
	if err != nil {
		return fmt.Errorf("Could not find DFX executable: %w", err)
	}

	dfxCommand := &exec.Cmd{
		Path:   dfxExecutable,
		Dir:    s.config.CanisterName,
		Args:   []string{dfxExecutable, "start", "--background"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := dfxCommand.Run(); err != nil {
		s.log.WithError(err).Errorln("Could not start local network")
		return err
	}
	s.log.Infoln("Sleeping 5 seconds to allow local network to set up enough that we can make requests...")
	time.Sleep(5 * time.Second)
	return nil
}

func (s *DFXService) createWriterIdentityIfNeeded() error {
	s.log.Infof("Creating writer identity...")
	output, exitCode, err := dfxCall(s.config.CanisterName, []string{"identity", "new", "writer"}, true)
	if err != nil {
		s.log.WithError(err).Errorln("Could not create new writer identity")
		return err
	}
	if exitCode != 0 && !strings.HasPrefix(output, "Creating identity: \"writer\".\nIdentity already exists.") {
		s.log.Errorln("Could not create new writer identity:", output)
		return err
	}
	return nil
}

func (s *DFXService) doesCanisterExist() (bool, error) {
	s.log.Infof("Checking if canister already exists...")
	output, exitCode, err := dfxCall(s.config.CanisterName, []string{"canister", "id", s.config.CanisterName}, true)
	if err != nil {
		s.log.WithError(err).Errorln("Could not determine if canister exists")
		return false, err
	}
	if exitCode == 0 {
		return true, nil
	}
	if !strings.HasPrefix(output, "Cannot find canister id.") {
		s.log.Errorln("Could not determine if canister exists:", output)
		return false, err
	}
	return false, nil
}

func (s *DFXService) isCanisterRunning() (bool, error) {
	s.log.Infof("Checking if canister is running...")
	output, _, err := dfxCall(s.config.CanisterName, []string{"canister", "status", s.config.CanisterName}, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not determine canister status:", output)
		return false, err
	}
	if !strings.HasPrefix(output, "Canister status call result for " + s.config.CanisterName) {
		s.log.WithError(err).Errorln("Could not determine canister status:", output)
		return false, fmt.Errorf("Could not determine canister status: %v", output)
	}
	isRunning := strings.HasPrefix(output, "Status: Running")
	return isRunning, nil
}

func (s *DFXService) createCanister() error {
	s.log.Infof("Creating canister...")
	output, _, err := dfxCall(s.config.CanisterName, []string{"canister", "create", s.config.CanisterName}, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not create canister:", output)
		return err
	}
	return nil
}

func (s *DFXService) buildCanister() error {
	s.log.Infof("Building canister...")
	output, _, err := dfxCall(s.config.CanisterName, []string{"build", s.config.CanisterName}, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not build canister:", output)
		return err
	}
	return nil
}

func (s *DFXService) installCanister(upgradeExistingCanister bool) error {
	s.log.Infof("Installing canister...")

	args := []string{"canister", "install", s.config.CanisterName}
	if upgradeExistingCanister {
		args = []string{"canister", "install", s.config.CanisterName, "--mode", "upgrade"}
	}

	output, _, err := dfxCall(s.config.CanisterName, args, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not install canister:", output)
		return err
	}
	return nil
}

func (s *DFXService) startCanister() error {
	s.log.Infof("Starting canister...")
	output, _, err := dfxCall(s.config.CanisterName, []string{"canister", "start", s.config.CanisterName}, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not start canister:", output)
		return err
	}
	return nil
}

func (s *DFXService) checkIsOwner() (bool, error) {
	s.log.Infof("Checking if we have the owner role...")
	output, _, err := dfxCall(s.config.CanisterName, []string{"canister", "call", s.config.CanisterName, "my_role"}, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not retrieve current role:", output)
		return false, err
	}
	hasOwnerRole := strings.HasPrefix(output, "(opt variant { owner })")
	return hasOwnerRole, nil
}

func (s *DFXService) assignOwnerRole() error {
	s.log.Infof("Assigning owner role to owner identity...")
	output, _, err := dfxCall(s.config.CanisterName, []string{"canister", "call", s.config.CanisterName, "assign_owner_role"}, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not assign owner role to the owner identity:", output)
		return err
	}
	return nil
}

func (s *DFXService) getWriterIDPrincipal() (string, error) {
	s.log.Infof("Retrieving writer ID principal...")
	output, _, err := dfxCall(s.config.CanisterName, []string{"--identity", "writer", "identity", "get-principal"}, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not determine the writer identity's principal:", output)
		return "", err
	}
	return strings.TrimSpace(output), nil
}

func (s *DFXService) assignWriterRole(writerPrincipal string) error {
	s.log.Infof("Assigning writer role to writer identity %s...", writerPrincipal)
	callArgs := fmt.Sprintf("(%v)", utils.CandidPrincipal(writerPrincipal))
	output, _, err := dfxCall(s.config.CanisterName, []string{"canister", "call", s.config.CanisterName, "assign_writer_role", callArgs}, false)
	if err != nil {
		s.log.WithError(err).Errorln("Could not assign writer role to the writer identity:", output)
		return err
	}
	return nil
}

func (s *DFXService) updateValueInCanister(key string, val map[string]float64) error {
	s.log.Infof("Updating value in canister...")

	for k, v := range val {
		callArgs := fmt.Sprintf("(%v,%v,%v)", utils.CandidText(key), utils.CandidText(k), utils.CandidFloat64(v))
		output, _, err := dfxCall(s.config.CanisterName, []string{"--identity", "writer", "canister", "call", s.config.CanisterName, "update_map_value", callArgs}, false)
		if err != nil {
			s.log.WithError(err).Errorln("Could not update key", key, "field", k, "value", val, "in canister:", output)
			return err
		}
	}
	return nil
}

func dfxCall(workingDir string, args []string, allowNonzeroExitCode bool) (string, int, error) {
	dfxExecutable, err := exec.LookPath("dfx")
	if err != nil {
		return "", 0, fmt.Errorf("Could not find DFX executable: %w", err)
	}

	dfxCommand := &exec.Cmd{
		Path: dfxExecutable,
		Dir:  workingDir,
		Args: append([]string{dfxExecutable}, args...),
	}
	output, err := dfxCommand.CombinedOutput()
	if err != nil {
		if allowNonzeroExitCode {
			if exitError, ok := err.(*exec.ExitError); ok {
				return string(output), exitError.ExitCode(), nil
			}
		}
		return string(output), 0, fmt.Errorf("Running DFX command failed: %w", err)
	}
	return string(output), 0, nil
}
