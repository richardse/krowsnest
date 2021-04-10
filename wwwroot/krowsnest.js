

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
