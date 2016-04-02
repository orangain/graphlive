var seriesStyles = [
  // See: https://www.nig.ac.jp/color/gen/#color
  {strokeStyle: 'rgb(230, 159, 0)', fillStyle: 'rgba(230, 159, 0, 0.3)', lineWidth: 2},
  {strokeStyle: 'rgb(86, 180, 233)', fillStyle: 'rgba(86, 180, 233, 0.3)', lineWidth: 2},
  {strokeStyle: 'rgb(0, 158, 115)', fillStyle: 'rgba(0, 158, 115, 0.3)', lineWidth: 2}
];

var series = null;
var messageElement = document.getElementById('message');
var smoothie = new SmoothieChart();
smoothie.streamTo(document.getElementById('graph'), 1000);

var wsProtocol = (location.protocol === 'https:') ? 'wss://' : 'ws://';
var wsUri = wsProtocol + location.host + '/ws';
console.log(wsUri);

showMessage('Connnecting...');
connect();

function connect() {
  var websocket = new WebSocket(wsUri);
  websocket.onmessage = function(evt) {
    var metrics = JSON.parse(evt.data);
    //console.log(metrics);
    if (series === null || series.length != metrics.metrics.length) {
      resetSeries(metrics.metrics.length);
    }

    var t = new Date();
    t.setTime(metrics.time);
    var message = t.toLocaleString();

    metrics.metrics.forEach(function(metric, i) {
      series[i].append(metrics.time, metric.value);
      message += '\n' + metric.label + ': ' + metric.value;
    });

    showMessage(message);
  };
  websocket.onerror = function(evt) {
    //showMessage('Failed to connect');
  };
  websocket.onopen = function(evt) {
    showMessage('Sucessfully connected.');
  };
  websocket.onclose = function(evt) {
    var sleep = 5;
    showMessage('Connection closed. Try reconnecting every ' + sleep + ' seconds.');

    setTimeout(function() {
      connect();
    }, sleep * 1000);
  };
}

function resetSeries(numSeries) {
  if (series !== null) {
    series.forEach(function(s) {
      smoothie.removeTimeSeries(s);
    });
  }

  series = [];
  for (var i = 0; i < numSeries; i++) {
    var s = new TimeSeries();
    series.push(s);
    smoothie.addTimeSeries(s, seriesStyles[i]);
  }
}

function showMessage(text) {
  messageElement.innerHTML = '';
  text.split(/\n/).forEach(function(line) {
    messageElement.appendChild(document.createTextNode(line));
    messageElement.appendChild(document.createElement('br'));
  });
}
