var monitorURL = decodeURIComponent(window.location.search.split('=')[1])

fetch('/getStatus', {
  method: 'POST',
  mode: 'cors',
  cache: 'no-cache',
  credentials: 'same-origin',
  headers: {
    'Content-Type': 'application/json; charset=utf-8'
  },
  redirect: 'follow',
  referrer: 'no-referrer',
  body: JSON.stringify({
    url: monitorURL
  })
})
  .then((response) => response.json())
  .then((data) => {
    var svgElem = document.getElementById('diagram')
    var svg = d3.select('svg')
    var width = svgElem.width.baseVal.value
    var height = svgElem.height.baseVal.value
    var g = svg.append('g').attr('transform', 'translate(100,0)')

    var tree = d3.tree().size([ height, width - 160 ])

    var stratify = d3.stratify()

    var normalizedData = []

    function addToNormalizedData (d, parentId) {
      const id = 'n' + normalizedData.length
      normalizedData.push({
        parentId,
        id,
        name: d.name,
        active: d.active,
        status: d.status,
        latency: d.latency
      })
      if (d.dependencies instanceof Array) {
        d.dependencies.forEach((dep) => addToNormalizedData(dep, id))
      }
    }
    addToNormalizedData(data)
    refreshFaultyNodeList(data)

    var root = stratify(normalizedData).sort(function (a, b) {
      return a.height - b.height || a.id.localeCompare(b.id)
    })

    var link = g.selectAll('.link').data(tree(root).links()).enter().append('g')

    link
      .append('path')
      .attr('class', (d) => {
        return 'link' + ' ' + (d.target.data.active ? 'link-active' : 'link-inactive')
      })
      .attr(
        'd',
        d3
          .linkHorizontal()
          .x(function (d) {
            return d.y
          })
          .y(function (d) {
            return d.x
          })
      )

    link
      .append('text')
      .attr('class', 'latency')
      .text((d) => {
        const ms = Math.round(d.target.data.latency / 1e5) / 10
        return ms + 'ms'
      })
      .attr('dy', -3)
      .attr('x', (d) => {
        return (d.source.y + d.target.y) / 2
      })
      .attr('y', (d) => {
        return (d.source.x + d.target.x) / 2
      })

    var node = g
      .selectAll('.node')
      .data(root.descendants())
      .enter()
      .append('g')
      .attr('class', function (d) {
        let className = d.data.active ? 'active' : 'inactive'
        return className + ' ' + 'node' + (d.children ? ' node--internal' : ' node--leaf')
      })
      .attr('transform', function (d) {
        return 'translate(' + d.y + ',' + d.x + ')'
      })

    node.append('circle').attr('r', 2.5)

    node
      .append('text')
      .attr('dy', 3)
      .attr('x', function (d) {
        return d.children ? -8 : 8
      })
      .style('text-anchor', function (d) {
        return d.children ? 'end' : 'start'
      })
      .text(function (d) {
        return d.data.name
      })
  })

const faultyNodeList = document.getElementById('faulty-node-list')

const newLi = (txt) => {
  const li = document.createElement('li')
  const txtNode = document.createTextNode(txt)
  li.appendChild(txtNode)
  return li
}

const refreshFaultyNodeList = (data, listNode) => {
  listNode = listNode || faultyNodeList
  while (listNode.firstChild) {
    listNode.removeChild(listNode.firstChild)
  }
  const faultyNodes = data.dependencies.filter((d) => !d.active)
  if (faultyNodes.length === 0) {
    listNode.appendChild(newLi('<none>'))
    return
  }
  faultyNodes.forEach((d) => {
    listNode.appendChild(newLi(d.name + ': ' + d.status))
    if (d.dependencies instanceof Array) {
      const newList = document.createElement('ul')
      listNode.appendChild(newList)
      refreshFaultyNodeList(d, newList)
    }
  })
}
