const go = new Go();

function setup() {
    let w = window.innerWidth;
    let h = window.innerHeight;

    noLoop();
    createCanvas(w, h);
}

function draw() {
    if (window.goDraw)
        window.goDraw();
}

let urlParams = new URLSearchParams(window.location.search);
let debugPort = urlParams.get("connectDebugger");

if (debugPort) {
    let script = document.createElement("script");
    script.onload = () => {
        if (!window.amphionDebugger)
            return;

        window.amphionDebugger.start(debugPort);
    };
    script.src = "amphiondebug.js";
    document.body.appendChild(script);
}

if (WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });
} else {
    fetch("app.wasm").then(response =>
        response.arrayBuffer()
    ).then(bytes =>
        WebAssembly.instantiate(bytes, go.importObject)
    ).then(result =>
        go.run(result.instance)
    );
}

function commencePanic(reason, msg) {
    document.querySelector("canvas").style.display = "none";
    document.querySelector(".panic").style.display = "block";
    document.querySelector("#panic-reason").innerHTML = reason;
    document.querySelector("#panic-message").innerHTML = msg;
}