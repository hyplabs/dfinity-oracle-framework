package framework

const CodeTemplate = `// Import Base Modules
import AssocList "mo:base/AssocList";
import Error "mo:base/Error";
import List "mo:base/List";
import Option "mo:base/Option";
import Text "mo:base/Text";

shared (msg) actor class() {
    // Define custom types
    public type Role = {
        #owner;
        #writer;
    };

    // Application State
    private stable var owner: ?Principal = null;
    private stable var roles: AssocList.AssocList<Principal, Role> = List.nil();
    private stable var map: AssocList.AssocList<Text, AssocList.AssocList<Text, Float>> = List.nil();
    private stable var destructed: Bool = false;

    // Favorite Cities Functions
    public shared ({caller}) func update_map_value(k: Text, p: Text, v: Float): async() {
        await require_role(caller, ?#writer);
        await require_undestructed();

        var sublist: ?AssocList.AssocList<Text, Float> = AssocList.find<Text, AssocList.AssocList<Text, Float>>(map, k, text_eq);
        var newSublist: AssocList.AssocList<Text, Float> = List.nil();
        var text: Text = "";

        if (Option.isSome(sublist) == true) {
            let flattenedSublist = Option.flatten(sublist);
            newSublist := AssocList.replace<Text, Float>(flattenedSublist, p, text_eq, ?v).0;
            map := AssocList.replace<Text, AssocList.AssocList<Text, Float>>(map, k, text_eq, ?newSublist).0;
        } else if (Option.isNull(sublist) == true) {
            newSublist := AssocList.replace<Text, Float>(newSublist, p, text_eq, ?v).0;
            map := AssocList.replace<Text, AssocList.AssocList<Text, Float>>(map, k, text_eq, ?newSublist).0;
        };
    };

    public func get_map_value(k: Text): async ?AssocList.AssocList<Text, Float> {
        await require_undestructed();
        return AssocList.find<Text, AssocList.AssocList<Text, Float>>(map, k, text_eq);
    };

    public func get_map_field_value(k: Text, p: Text): async ?Float {
        await require_undestructed();
        let sublist = AssocList.find<Text, AssocList.AssocList<Text, Float>>(map, k, text_eq);
        if (Option.isSome(sublist)) {
            let flattenedSublist = Option.flatten(sublist);
            return AssocList.find<Text, Float>(flattenedSublist, p, text_eq);
        };
        return null;
    };

    public func get_map(): async AssocList.AssocList<Text, AssocList.AssocList<Text, Float>> {
        await require_undestructed();
        return map;
    };

    func text_eq(a: Text, b: Text): Bool {
        return a == b;
    };

    // Identity Access Control Functions
    func principal_eq(a: Principal, b: Principal): Bool {
        return a == b;
    };

    func get_role(p: Principal): ?Role {
        if (Option.isNull(owner) == true) {
            return null;
        } else if (?p == owner) {
            return ?#owner;
        } else {
            return AssocList.find<Principal, Role>(roles, p, principal_eq);
        };
    };

    func require_role(p: Principal, r: ?Role): async() {
        if(r != get_role(p)) {
            throw Error.reject("You do not have the required role to perform this operation.")
        };
    };

    func require_undestructed(): async() {
        if (destructed == true) {
            throw Error.reject("This oracle canister was destructed by the owner. It may have been corrupted or become malicious.")
        };
    };

    public shared ({caller}) func assign_owner_role(): async() {
        if (Option.isSome(owner) == true) {
            throw Error.reject("Cannot set owner if there is already an owner");
        };
        owner := ?caller;
    };

    public shared ({caller}) func assign_writer_role(p: Principal): async() {
        await require_role(caller, ?#owner);
        await require_undestructed();
        if (?p == owner) {
            throw Error.reject("Specified principal is the canister owner, which cannot also be the canister writer");
        };
        roles := AssocList.replace<Principal, Role>(roles, p, principal_eq, ?#writer).0;
    };

    public shared ({caller}) func revoke_writer_role(p: Principal): async() {
        await require_role(caller, ?#owner);
        await require_undestructed();
        if (?p == owner) {
            throw Error.reject("Specified principal is the canister owner, which cannot also be the canister writer");
        };
        if (get_role(p) != ?#writer) {
            throw Error.reject("Specified principal was not a writer to begin with");
        };
        roles := AssocList.replace(roles, p, principal_eq, null).0;
    };

    public shared ({caller}) func my_role(): async ?Role {
        return get_role(caller);
    };

    public shared ({caller}) func get_roles(): async List.List<(Principal, Role)> {
        await require_role(caller, ?#owner);
        await require_undestructed();
        return roles;
    };

    public shared ({caller}) func self_destruct(): async() {
        await require_role(caller, ?#owner);
        destructed := true;
        map := List.nil();
        roles := List.nil();
    }
}`
