let bingoData = null;
let currentSeed = null;
let checkedCells = new Set();

// FNV-1a hash (32-bit)
function fnv1a(str) {
    let hash = 2166136261;
    for (let i = 0; i < str.length; i++) {
        hash ^= str.charCodeAt(i);
        hash = Math.imul(hash, 16777619);
    }
    return hash >>> 0;
}

// Seeded random number generator (mulberry32)
function mulberry32(seed) {
    return function() {
        let t = seed += 0x6D2B79F5;
        t = Math.imul(t ^ t >>> 15, t | 1);
        t ^= t + Math.imul(t ^ t >>> 7, t | 61);
        return ((t ^ t >>> 14) >>> 0) / 4294967296;
    };
}

// Convert seed string to number
function seedToNumber(seed) {
    if (!seed || seed.length === 0) {
        return Date.now();
    }
    const num = parseInt(seed, 10);
    if (!isNaN(num)) {
        return num;
    }
    return fnv1a(seed);
}

// Generate 25 unique random indices
function generateIndices(rng, max) {
    const indices = [];
    const used = new Set();

    while (indices.length < 25) {
        const index = Math.floor(rng() * max);
        if (!used.has(index)) {
            used.add(index);
            indices.push(index);
        }
    }

    return indices;
}

// Generate bingo card
function generateBingoCard(seed) {
    if (!bingoData || bingoData.goals.length < 25) {
        console.error('Bingo data not loaded or insufficient goals');
        return null;
    }

    const seedNum = seedToNumber(seed);
    currentSeed = seed || seedNum.toString();

    const rng = mulberry32(seedNum);
    const indices = generateIndices(rng, bingoData.goals.length);

    return indices.map(i => bingoData.goals[i]);
}

// Render bingo card
function renderCard(goals) {
    const card = document.getElementById('bingo-card');
    card.innerHTML = '';

    goals.forEach((goal, index) => {
        const cell = document.createElement('div');
        cell.className = 'bingo-cell';
        cell.textContent = goal.name;
        cell.dataset.index = index;

        if (checkedCells.has(index)) {
            cell.classList.add('checked');
        }

        cell.addEventListener('click', () => toggleCell(cell, index));
        card.appendChild(cell);
    });

    document.getElementById('current-seed').textContent = currentSeed;
}

// Toggle cell checked state
function toggleCell(cell, index) {
    if (checkedCells.has(index)) {
        checkedCells.delete(index);
        cell.classList.remove('checked');
    } else {
        checkedCells.add(index);
        cell.classList.add('checked');
    }
}

// Generate and render new card
function newCard(seed) {
    checkedCells.clear();
    const goals = generateBingoCard(seed);
    if (goals) {
        renderCard(goals);
        updateURL();
    }
}

// Update URL with current seed
function updateURL() {
    const url = new URL(window.location);
    url.searchParams.set('seed', currentSeed);
    history.replaceState(null, '', url);
}

// Copy share URL
function copyShareURL() {
    const url = new URL(window.location);
    url.searchParams.set('seed', currentSeed);
    navigator.clipboard.writeText(url.toString()).then(() => {
        alert('URLをコピーしました');
    });
}

// Initialize
async function init() {
    try {
        const response = await fetch('bingo.json');
        bingoData = await response.json();

        const params = new URLSearchParams(window.location.search);
        const seed = params.get('seed') || '';

        document.getElementById('seed-input').value = seed;
        newCard(seed);

        document.getElementById('generate-btn').addEventListener('click', () => {
            const seed = document.getElementById('seed-input').value;
            newCard(seed);
        });

        document.getElementById('share-btn').addEventListener('click', copyShareURL);

        document.getElementById('seed-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                newCard(e.target.value);
            }
        });
    } catch (error) {
        console.error('Failed to load bingo data:', error);
        document.getElementById('bingo-card').innerHTML = '<p>データの読み込みに失敗しました</p>';
    }
}

document.addEventListener('DOMContentLoaded', init);
