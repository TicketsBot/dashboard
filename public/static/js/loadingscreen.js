function showLoadingScreen() {
    const content = document.getElementsByClassName('content')[0];
    content.style.display = 'none';
    document.getElementById('loading-container').style.display = 'block';
}

function hideLoadingScreen() {
    document.getElementById('loading-container').style.display = 'none';

    const content = document.getElementsByClassName('content')[0];
    content.style.display = 'block';
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
