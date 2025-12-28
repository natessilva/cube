const worker = new Worker("worker.js");

worker.onmessage = (event) => {
  solution.innerText = event.data.solution;
  duration.innerText = event.data.duration;
  const solutionText = event.data.solution || '';
  const moveCount = solutionText.trim().split(/\s+/).filter(move => move !== '').length;
  moves.innerText = `Moves: ${moveCount}`;
};

worker.onerror = (error) => {
  console.error("Worker error:", error);
};

const scramble = document.getElementById("scramble");
scramble.value = generateScramble();
const solution = document.getElementById("solution");
const duration = document.getElementById("duration");
const moves = document.getElementById("moves");
const form = document.getElementById("form");
const randomize = document.getElementById("randomize");

form.addEventListener("submit", (e) => {
  worker.postMessage(scramble.value);
  solution.innerText = "Solving...";
  e.preventDefault();
});

function generateScramble() {
  const faces = ["U", "D", "L", "R", "F", "B"];
  const modifiers = ["", "'", "2"]; // No modifier, prime, or double turn
  const moves = Math.floor(Math.random() * 8) + 23; // Random number between 23 and 30
  const scramble = [];

  for (let i = 0; i < moves; i++) {
    let nextMove;
    do {
      const face = faces[Math.floor(Math.random() * faces.length)];
      const modifier = modifiers[Math.floor(Math.random() * modifiers.length)];
      nextMove = face + modifier;
    } while (i > 0 && scramble[i - 1][0] === nextMove[0]); // Avoid consecutive moves on the same face

    scramble.push(nextMove);
  }

  return scramble.join(" ");
}

randomize.addEventListener("click", () => {
  scramble.value = generateScramble();
});
