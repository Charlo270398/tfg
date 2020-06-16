const tagsArray = [];
const tagsArrayNumeros = [];
var elementoSeleccionado = -1;

function init () {
  deleteBreadcrumb();
  addLinkBreadcrumb('Usuario', '/user/menu');
  addLinkBreadcrumb('Medico', '/user/doctor');
  addLinkBreadcrumb('Investigación', '');
  const elementoSelector = document.querySelector("#elementoSelector");
  elementoSelector.addEventListener('change',cambioElemento,false);
}

document.addEventListener('DOMContentLoaded',init,false);

function cambioElemento(event){
  elementoSeleccionado = event.target.value;
  if(tagsArray.length > 0){
    changeGraphics();
  }
}

var sample = [];

var svg = d3.select('svg');
var svgContainer = d3.select('#container');

var margin = 80;
var width = 1000 - 2 * margin;
var height = 600 - 2 * margin;

var chart = svg.append('g')
  .attr('transform', `translate(${margin}, ${margin})`);

var xScale = d3.scaleBand()
  .range([0, width])
  .domain(sample.map((s) => s.elemento))
  .padding(0.4)

  var yScale = d3.scaleLinear()
  .range([height, 0])
  .domain([0, 20]);

// vertical grid lines
// const makeXLines = () => d3.axisBottom()
//   .scale(xScale)

var makeYLines = () => d3.axisLeft()
  .scale(yScale)

function loadGraphics(nombreElemento){
    chart.append('g')
    .attr('transform', `translate(0, ${height})`)
    .call(d3.axisBottom(xScale));

  chart.append('g')
    .call(d3.axisLeft(yScale));


  chart.append('g')
    .attr('class', 'grid')
    .call(makeYLines()
      .tickSize(-width, 0, 0)
      .tickFormat('')
    )

  var barGroups = chart.selectAll()
    .data(sample)
    .enter()
    .append('g')

  barGroups
    .append('rect')
    .attr('class', 'bar')
    .attr('x', (g) => xScale(g.elemento))
    .attr('y', (g) => yScale(g.value))
    .attr('height', (g) => height - yScale(g.value))
    .attr('width', xScale.bandwidth())
    .on('mouseenter', function (actual, i) {
      d3.selectAll('.value')
        .attr('opacity', 0)

      d3.select(this)
        .transition()
        .duration(300)
        .attr('opacity', 0.6)
        .attr('x', (a) => xScale(a.elemento) - 5)
        .attr('width', xScale.bandwidth() + 10)

      const y = yScale(actual.value)

      line = chart.append('line')
        .attr('id', 'limit')
        .attr('x1', 0)
        .attr('y1', y)
        .attr('x2', width)
        .attr('y2', y)

      barGroups.append('text')
        .attr('class', 'divergence')
        .attr('x', (a) => xScale(a.elemento) + xScale.bandwidth() / 2)
        .attr('y', (a) => yScale(a.value) + 30)
        .attr('fill', 'white')
        .attr('text-anchor', 'middle')
        .text((a, idx) => {
          const divergence = (a.value - actual.value).toFixed(1)
          
          let text = ''
          if (divergence > 0) text += '+'
          text += `${divergence}`

          return idx !== i ? text : '';
        })

    })
    .on('mouseleave', function () {
      d3.selectAll('.value')
        .attr('opacity', 1)

      d3.select(this)
        .transition()
        .duration(300)
        .attr('opacity', 1)
        .attr('x', (a) => xScale(a.elemento))
        .attr('width', xScale.bandwidth())

      chart.selectAll('#limit').remove()
      chart.selectAll('.divergence').remove()
    })

  barGroups 
    .append('text')
    .attr('class', 'value')
    .attr('x', (a) => xScale(a.elemento) + xScale.bandwidth() / 2)
    .attr('y', (a) => yScale(a.value) + 0)
    .attr('text-anchor', 'middle')
    .text((a) => `${a.value}`)
  
  svg
    .append('text')
    .attr('class', 'label')
    .attr('x', -(height / 2) - margin)
    .attr('y', margin / 2.4)
    .attr('transform', 'rotate(-90)')
    .attr('text-anchor', 'middle')
    .text('Valores de ' + nombreElemento + " (g/dL)")

  svg.append('text')
    .attr('class', 'label')
    .attr('x', width / 2 + margin)
    .attr('y', height + margin * 1.7)
    .attr('text-anchor', 'middle')
    .text('Valores')

  svg.append('text')
    .attr('class', 'title')
    .attr('x', width / 2 + margin)
    .attr('y', 40)
    .attr('text-anchor', 'middle')
    .text(nombreElemento)
  }
  
  function changeCheckbox(event){
    if(event.target.checked){
      tagsArrayNumeros.push(parseInt(event.target.parentNode.parentNode.getAttribute('tag_id')));
      tagsArray.push(event.target.parentNode.parentNode.textContent);
    }else{
        var pos = tagsArrayNumeros.indexOf(parseInt(event.target.parentNode.parentNode.getAttribute('tag_id')));
        tagsArrayNumeros.splice(pos, 1);
        pos = tagsArray.indexOf(event.target.parentNode.parentNode.textContent);
        tagsArray.splice(pos, 1);
    }

    if(elementoSeleccionado != -1){
      changeGraphics();
    }
  }

  function changeGraphics(){
    while (document.querySelector("#svg").lastElementChild) { 
        document.querySelector("#svg").removeChild(document.querySelector("#svg").lastElementChild); 
    } 
    //Modificamos datos a seleccionar
    sample = [];
    tagsArrayNumeros.forEach(tag => {
        var valor = 0.0;
        var contador = 0;
        analiticas.forEach(analitica => {
            if(analitica.tags.includes(tag)){
                switch (elementoSeleccionado) {
                    case "1":
                        valor += parseFloat(analitica.leucocitos);
                        contador++;
                        break;
                    case "2":
                        valor += parseFloat(analitica.hematies);
                        contador++;
                        break;
                    case "3":
                        valor += parseFloat(analitica.plaquetas);
                        contador++;
                        break;
                    case "4":
                        valor += parseFloat(analitica.glucosa);
                        contador++;
                        break;
                    case "5":
                        valor += parseFloat(analitica.hierro);
                        contador++;
                        break;
                    default:
                        break;
                }
            }
        });
        if(contador != 0){
            sample.push({elemento: tagsArray[tagsArrayNumeros.indexOf(tag)], value: (valor/contador)});
        }else{
            sample.push({elemento: tagsArray[tagsArrayNumeros.indexOf(tag)], value: 0.0});
        }
    });
  
    if(elementoSeleccionado != -1){
        svg = d3.select('svg');
        svgContainer = d3.select('#container');
        
        margin = 80;
        width = 1000 - 2 * margin;
        height = 600 - 2 * margin;

        chart = svg.append('g')
            .attr('transform', `translate(${margin}, ${margin})`);

        xScale = d3.scaleBand()
            .range([0, width])
            .domain(sample.map((s) => s.elemento))
            .padding(0.4)
        
        yScale = d3.scaleLinear()
            .range([height, 0])
            .domain([0, 20]);

        makeYLines = () => d3.axisLeft()
            .scale(yScale)
        loadGraphics(getNombreElemento(elementoSeleccionado));
    }
}

function getNombreElemento(elementoNumero){
  switch (elementoNumero) {
      case "1":
          return "Leucocitos";
      case "2":
          return "Hematies";
      case "3":
          return "Plaquetas";
      case "4":
          return "Glucosa";
      case "5":
          return "Hierro";
      default:
          return "NO HAY ELEMENTO"
  }
}