function showLoadingScreen() {
    const content = document.getElementsByClassName('content')[0] || document.getElementsByClassName('tcontent-container')[0];
    content.style.display = 'none';
    document.getElementById('loading-container').style.display = 'block';
}

function hideLoadingScreen() {
    document.getElementById('loading-container').style.display = 'none';

    const content = document.getElementsByClassName('content')[0] || document.getElementsByClassName('tcontent-container')[0];
    if (content.classList.contains('tcontent-container')) {
        content.style.display = 'flex';
    } else {
        content.style.display = 'block';
    }

    content.classList.add('fade-in');
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function withLoadingScreen(func) {
    showLoadingScreen();
    await func();
    hideLoadingScreen();
}
