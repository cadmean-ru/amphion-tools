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