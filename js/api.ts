declare var require: any
declare var Promise: any;
declare var $: any;
declare var global: any;
declare var Gremlin: any;
declare var request: any;
var makeRequest = null;

import jsEnv = require('browser-or-node');
if (jsEnv) {
    // When using node, we use najax for $.ajax calls
    if (jsEnv.isNode && !jsEnv.isBrowser) {
        var najax = require('najax');
        global.$ = <any>{ "ajax": najax };
    }

    makeRequest = function (url: string, method: string, data: any) : any {
        return new Promise(function (resolve, reject) {
            return $.ajax({
                dataType: "json",
                url: "http://localhost:8082" + url,
                data: data,
                contentType: "application/json; charset=utf-8",
                method: method
            })
            .done(function (data, statusText, jqXHR) {
                if (jqXHR.status == 200) {
                    resolve(data);
                } else {
                    reject(data);
                }
            });
        })
    }
} else {
    // We're running inside Otto. We use the 'request' method
    // exported by Skydive that makes use of the REST client
    // so we get support for authentication
    makeRequest = function (url: string, method: string, data: string) : any {
        return new Promise(function (resolve, reject) {
            var err: any
            var output: any
            try {
                output = request(url, method, data)
                if (output && output.length > 0) {
                    output = JSON.parse(output)
                } else {
                    output = null
                }
                resolve(output)
            }
            catch(e) {
                reject(e)
            }
        })
    }
}

function paramsString(params: any[]): string {
  var s = ""
  for (let param in params) {
      if (param != "0") {
          s += ", "
      }
      switch (typeof params[param]) {
          case "string":
              s += '"' + params[param] + '"';
              break;
          default:
              s += params[param];
              break;
      }
  }
  return s
}
// Returns a string of the invocation of 'method' with its parameters 'params'
// like MyMethod(1, 2, "s")
function callString(method: string, ...params: any[]): string {
    let s = method + "("
    s += paramsString(params)
    s += ")"
    return s
}

export class SerializationHelper {
    static toInstance<T>(obj: T, jsonObj: Object) : T {
        if (typeof obj["fromJSON"] === "function") {
            obj["fromJSON"](jsonObj);
        }
        else {
            for (var propName in jsonObj) {
                obj[propName] = jsonObj[propName]
            }
        }

        return obj;
    }

    static unmarshalArray<T>(data, c: new () => T): T[] {
        var items: T[] = [];
        for (var obj in data) {
            let newT: T = SerializationHelper.toInstance(new c(), data[obj]);
            items.push(newT);
        }
        return items
    }

    static unmarshalMap<T>(data, c: new () => T): { [key: string]: T; } {
        var items: { [key: string]: T; };
        for (var obj in data) {
            let newT: T = SerializationHelper.toInstance(new c(), data[obj]);
            items[obj] = newT;
        }
        return items
    }
}

export interface APIObject {
    UUID: string;
}

export class Alert implements APIObject {
    UUID: string;
    Name: string;
    Description: string;
    Expression:  string;
    Action: string;
    Trigger: string;
}

export class Capture implements APIObject {
    UUID: string
    GremlinQuery: string
    BPFFilter: string
    Name: string
    Description: string
    Type: string
    Count: string
    PCAPSocket: string
    Port: string
    RawPacketLimit: string
    HeaderSize: string
    ExtraTCPMetric: string
    IPDefrag: boolean
    ReassembleTCP: boolean
}

export class PacketInjection implements APIObject {
    UUID:       string
    Src:        string
    Dst:        string
    SrcIP:      string
    DstIP:      string
    SrcMAC:     string
    DstMAC:     string
    SrcPort:    Number
    DstPort:    Number
    Type:       string
    Payload:    string
    TrackingID: string
    ICMPID:     Number
    Count:      Number
    Interval:   Number
    Increment:  boolean
    StartTime:  Date
}

export class API<T extends APIObject> {
    Resource: string
    Factory: new () => T

    constructor(resource: string, c: new () => T) {
        this.Resource = resource;
        this.Factory = c;
    }

    create(obj: T) {
        let resource = this.Resource;
        return makeRequest("/api/" + resource, "POST", JSON.stringify(obj))
            .then(function (data) {
                return SerializationHelper.toInstance(new (<any>obj.constructor)(), data);
            })
    }

    list() {
        let resource = this.Resource;
        let factory = this.Factory;

        return makeRequest('/api/' + resource, "GET", "")
            .then(function (data) {
                let resources = { };
                for (var obj in data) {
                    let resource = SerializationHelper.toInstance(new factory(), data[obj]);
                    resources[obj] = resource
                }
                return resources
            });
    }

    get(id) {
        let resource = this.Resource;
        let factory = this.Factory;

        return makeRequest('/api/' + resource + "/" + id, "GET", "")
            .then(function (data) {
                return SerializationHelper.toInstance(new factory(), data)
            })
            .catch(function (error) {
                return Error("Capture not found")
            })
    }

    delete(UUID: string) {
        let resource = this.Resource;
        return makeRequest('/api/' + resource + "/" + UUID, "DELETE", "")
    }
}

export class GraphElement {
    ID: string
    Host: string
    Metadata: Object
    CreatedAt: Date
    UpdatedAt: Date
    Revision: number
}

export class GraphNode extends GraphElement {
}

export class GraphEdge extends GraphElement {
    Parent: string
    Child: string
}

export class Graph {
    Nodes: GraphNode[]
    Edges: GraphEdge[]
}

export class GremlinAPI {
    query(s: string) {
        return makeRequest('/api/topology', "POST", JSON.stringify({'GremlinQuery': s}))
            .then(function (data) {
                if (data === null)
                    return [];
                // Result can be [Node] or [[Node, Node]]
                if (data.length > 0 && data[0] instanceof Array)
                    data = data[0];
                return data
            });
    }

    G() : G {
        return new G(this)
    }
}

export class _Metadata {
    params: any[]

    constructor(params: any[]) {
      this.params = params
    }

    public toString(): string {
        return "Metadata(" + paramsString(this.params) + ")"
    }
}

export function Metadata(params: any[]) : _Metadata {
  return new _Metadata(params)
}

export interface Step {
    name(): string
    serialize(data): any
}

export class Step implements Step {
    api: GremlinAPI
    gremlin: string
    previous: Step
    params: any[]

    name() { return "InvalidStep"; }

    constructor(api: GremlinAPI, previous: Step = null, ...params: any[]) {
        this.api = api
        this.previous = previous
        this.params = params
    }

    public toString(): string {
        if (this.previous !== null) {
            let gremlin = this.previous.toString() + ".";
            return gremlin + callString(this.name(), ...this.params)
        } else {
            return this.name()
        }
    }

    then(cb) {
        let self = this;
        return self.api.query(self.toString())
            .then(function (data) {
                cb(data)
            });
    }
}

export class G extends Step {
    name() { return "G" }

    constructor(api: GremlinAPI) {
        super(api)
    }

    serialize(data) {
        return {
            "Nodes": SerializationHelper.unmarshalArray(data[0].Nodes, GraphNode),
            "Edges": SerializationHelper.unmarshalArray(data[0].Edges, GraphEdge),
        };
    }

    V(...params: any[]): V {
        return new V(this.api, this, ...params);
    }

    E(...params: any[]): E {
        return new E(this.api, this, ...params);
    }

    Flows(...params: any[]): Flows {
        return new Flows(this.api, this, ...params);
    }
}

export class V extends Step {
    name() { return "V" }

    serialize(data) {
        return SerializationHelper.unmarshalArray(data, GraphNode);
    }

    Has(...params: any[]): Has {
        return new Has(this.api, this, ...params);
    }

    ShortestPathTo(...params: any[]): ShortestPath {
        return new ShortestPath(this.api, this, ...params);
    }
}

export class E extends Step {
    name() { return "E" }

    serialize(data) {
        return SerializationHelper.unmarshalArray(data, GraphEdge);
    }

    Has(...params: any[]): Has {
        return new Has(this.api, this, ...params);
    }
}

export class Has extends V {
    name() { return "Has" }

    serialize(data) {
        return this.previous.serialize(data);
    }
}

export class ShortestPath extends Step {
    name() { return "ShortestPathTo" }

    public toString(): string {
        if (this.previous !== null) {
            let gremlin = this.previous.toString() + ".";
            return gremlin + callString(this.name(), ...this.params)
        } else {
            return this.name()
        }
    }

    serialize(data) {
        var items: GraphNode[][] = [];
        for (var obj in data) {
            items.push(SerializationHelper.unmarshalArray(data, GraphNode));
        }
        return items;
    }
}

class Predicate {
    name: string
    params: any

    constructor(name: string, ...params: any[]) {
        this.name = name
        this.params = params
    }

    public toString(): string {
        return callString(this.name, ...this.params)
    }
}

export function NE(param: any): Predicate {
    return new Predicate("NE", param)
}

export function GT(param: any): Predicate {
    return new Predicate("GT", param)
}

export function LT(param: any): Predicate {
    return new Predicate("LT", param)
}

export function GTE(param: any): Predicate {
    return new Predicate("GTE", param)
}

export function LTE(param: any): Predicate {
    return new Predicate("LTE", param)
}

export function IPV4(param: any): Predicate {
    return new Predicate("IPV4", param)
}

export function REGEX(param: any): Predicate {
    return new Predicate("REGEX", param)
}

export function NULL(): Predicate {
    return new Predicate("NULL")
}

export class FlowLayer {
    Protocol: string
    A: string
    B: string
    ID: string
}

export class ICMPLayer {
    Type: string
    Code: number
    ID: number
}

export class FlowMetric {
    ABPackets: number
    ABBytes: number
    BAPackets: number
    BABytes: number
    Start: number
    Last: number
}

export class RawPacket {
    Timestamp: number
    Index: number
}

export class TCPMetric {
    ABSynStart: number
    BASynStart: number
    ABSynTTL: number
    BASynTTL: number
    ABFinStart: number
    BAFinStart: number
    ABRstStart: number
    BARstStart: number
}

export class Flow {
    UUID: string
    LayersPath: string
    Application: string

    Link: FlowLayer
    Network: FlowLayer
    Transport: FlowLayer
    ICMP: ICMPLayer

    LastUpdateMetric: FlowMetric
    Metric: FlowMetric

    TCPFlowMetric: TCPMetric

    Start: number
    Last: number
    RTT: number
    TrackingID: string
    L3TrackingID: string
    ParentUUID: string

    NodeTID: string
    ANodeTID: string
    BNodeTID: string

    LastRawPackets: RawPacket[]

    RawPacketsCaptured: number
}

export class Flows extends Step {
    name() { return "Flows" }

    serialize(data) {
        return SerializationHelper.unmarshalArray(data, Flow);
    }

    Has(...params: any[]): Has {
        return new Has(this.api, this, ...params);
    }
}
