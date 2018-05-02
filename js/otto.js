var gremlinAPI = new GremlinAPI()
var G = gremlinAPI.G()

function objectToArray(obj) {
  var array = []
  for (var key in obj) {
    array.push(key)
    array.push(obj[key])
  }
  return array
}

function checkMTU(fromMetadata, toMetadata) {
  var from
  var MTU

  var V = G.V()
  from = V.Has.apply(V, objectToArray(fromMetadata))

  var paths = from.ShortestPathTo(Metadata(objectToArray(toMetadata))).then(function (path) {
    for (var i in path) {
      var node = path[i]
      if (MTU != undefined && (node.Metadata == undefined || node.Metadata.MTU > MTU)) {
        console.log("MTU " + node.Metadata.MTU + " on node " + node.ID + " is inferior to " + MTU)
        return false
      }
      MTU = node.Metadata.MTU
    }
  })

  return true
}
