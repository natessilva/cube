importScripts("wasm_exec.js");

const go = new Go();
let instance = null;

async function initializeWasm() {
  if (!instance) {
    try {
      const wasmModule = await WebAssembly.instantiateStreaming(
        fetch("main.wasm"),
        go.importObject
      );
      instance = wasmModule.instance;
      go.run(instance);
    } catch (err) {
      console.error("Failed to initialize Wasm module:", err);
    }
  }
}

onmessage = async (event) => {
  const message = event.data;

  if (typeof message === "string") {
    if (!instance) {
      await initializeWasm();
    }

    try {
      const result = globalThis.solve(message);
      postMessage(result);
    } catch (err) {
      console.error("Error calling Go function:", err);
    }
  } else {
    console.error("Received non-string message:", message);
  }
};

initializeWasm();
