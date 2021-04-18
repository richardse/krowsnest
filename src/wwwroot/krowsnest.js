function init() {
  if (!namespace_url) {
    alert('Namespace_url not defined. Make sure that config.js is loaded first.');
    throw new Error;
  }
  initNsPicker();
  initMermaid();
}

function initNsPicker() {
  fetch(namespace_url)
    .then(response => response.json())
    .then(nslist => {
      const nsSelect = document.getElementById('selectns');
      while (nsSelect.options.length) nsSelect.remove(0);

      Object.values(nslist).forEach(ns => {
        var opt = document.createElement('option');
        opt.appendChild(document.createTextNode(ns));
        opt.value = ns;
        nsSelect.appendChild(opt); 
      });
    });
}

function initMermaid() {
  var insertSvg = function(svgCode, bindFunctions) {
    const graphElement = document.getElementById("graph");
    graphElement.innerHTML = svgCode;
    bindFunctions(graphElement);
  }
  
  var config = {
    startOnLoad: false,
    securityLevel: "antiscript",
    flowchart:{
        useMaxWidth: false,
        htmlLabels: true
    }
  };
  mermaid.initialize(config);
  
  fetch('/mermaid')
    .then(response => response.text())
    .then(data => {
      mermaid.mermaidAPI.render("graphDiv", data, insertSvg);
    });
}
