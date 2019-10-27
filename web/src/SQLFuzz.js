import React, { useEffect } from 'react';
import ReactDOM from 'react-dom';
import G6 from '@antv/g6';
import * as d3 from 'd3';
import './sql_fuzz.css'

function SQLFuzz() {
  const ref = React.useRef(null)

  useEffect(() => {
    function ellipseContent(content) {
      if (content.length > 6) {
        return content.slice(0, 4) + '...'
      }
      return content
    }

    function fetchData(url) {
      return fetch(url).then(res => res.json())
    }

    Promise.all([
      fetchData('http://localhost:43000/graph'),
      fetchData('http://localhost:43000/heat')
    ]).then(([graphs, heatsArr]) => {
      // d3.select('#mountNode').text('')

      // https://bl.ocks.org/pstuffa/d5934843ee3a7d2cc8406de64e6e4ea5
      const heatValuesSet = heatsArr.map(h => h.heat)
      const minHeat = d3.min(heatValuesSet)
      const maxHeat = d3.max(heatValuesSet)
      const colorScale = d3
          .scaleSequential(d3.interpolateOranges)
          .domain([minHeat, maxHeat])

      const heatValuesMap = heatsArr.reduce((acc, cur) => {
        acc[`${cur.number}_${cur.alter}`] = { heat: cur.heat, sql: cur.sql }
        return acc
      }, {})

      let data = { nodes: [], edges: [] }

      /*node*/
      if (graphs.length !== 0) {
        let x = 100,
            y = window.innerHeight / 2
        for (let item of graphs) {
          let node = {}

          node.id = `${item.number}`
          node.description = 'label: ' + item.head
          node.label = ellipseContent(item.head)
          node.shape = 'circle'
          node.x = x
          node.y = y
          data.nodes.push(node)
          data.nodes.push(node)

          x = x + 100
          y = 50

          if (item.alter && item.alter.length) {
            for (let i = 0; i < item.alter.length; i++) {
              let subNode = {}
              subNode.id = `${node.id}_${i}`
              subNode.description = 'label: ' + item.alter[i].content
              subNode.label = ellipseContent(item.alter[i].content)
              subNode.x = x
              subNode.y = y
              data.nodes.push(subNode)
              if(y > window.innerHeight / 2 - 100 && y < window.innerHeight / 2) {
                y = window.innerHeight / 2 + (window.innerHeight / 2 - y)
              }else {
                y += 50
              }
              data.edges.push({ source: node.id, target: subNode.id })

              let fanout = item.alter[i].fanout
              for (let item of fanout) {
                data.edges.push({ source: subNode.id, target: `${item}` })
              }
            }
            y = window.innerHeight / 2
          }
        }
      }

      /*heat*/
      if (heatsArr.length !== 0) {
        data.nodes = data.nodes.map(node => {
          if (heatValuesMap[node.id] !== undefined) {
            const heatValue = heatValuesMap[node.id]
            return {
              ...node,
              description: `${node.description}<br/>heat: ${heatValue.heat}`,
              style: {
                lineWidth: 2,
                fill: colorScale(heatValue.heat)
              },
            }
          }
          return node
        })
      }

      /*pic*/
      const graph = new G6.Graph({
        container: ReactDOM.findDOMNode(ref.current),
        width: window.innerWidth,
        height: window.innerHeight,

        defaultNode: {
          size: [40, 20],
          shape: 'rect'
        },
        defaultEdge: {
          size: 2,
          color: '#e2e2e2',
          style: {
            endArrow: {
              path: 'M 4,0 L -4,-4 L -4,4 Z',
              d: 4
            }
          }
        },
        edgeStateStyles: {
          highlight: {
            stroke: '#999'
          }
        },
        modes: {
          default: [
            'drag-node',
            'drag-canvas',
            'zoom-canvas',
            {
              type: 'tooltip',
              formatText: function formatText(model) {
                return model.description
              },
              shouldUpdate: function shouldUpdate(e) {
                return true
              }
            }
          ]
        }
      })

      /* Display associated nodes */
      function clearAllStats() {
        graph.setAutoPaint(false);
        graph.getNodes().forEach(function(node) {
          graph.clearItemStates(node);
        });
        graph.getEdges().forEach(function(edge) {
          graph.clearItemStates(edge);
        });
        graph.paint();
        graph.setAutoPaint(true);
      }
      graph.on('node:mouseenter', function(e) {
        const item = e.item;
        graph.setAutoPaint(false);
        graph.getEdges().forEach(function(edge) {
          if (edge.getSource() === item) {
            graph.setItemState(edge.getTarget(), 'dark', false);
            graph.setItemState(edge.getTarget(), 'highlight', true);
            graph.setItemState(edge, 'highlight', true);
            edge.toFront();
          } else if (edge.getTarget() === item) {
            graph.setItemState(edge.getSource(), 'dark', false);
            graph.setItemState(edge.getSource(), 'highlight', true);
            graph.setItemState(edge, 'highlight', true);
            edge.toFront();
          } else {
            graph.setItemState(edge, 'highlight', false);
          }
        });
        graph.paint();
        graph.setAutoPaint(true);
      });
      graph.on('node:mouseleave', clearAllStats);
      graph.data(data)
      graph.render()
    })

  }, [])

  return (
      <div>
        <div ref={ref}></div>
      </div>
  );
}
export default SQLFuzz