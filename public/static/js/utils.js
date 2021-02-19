function showToast(title, content) {
    const container = document.getElementById('toast-container');

    container.innerHTML += `
            <div class="toast" role="alert" aria-live="assertive" aria-atomic="true" data-autohide="false">
              <div class="toast-header">
                <strong class="mr-auto">${title}</strong>
                <button type="button" class="ml-2 mb-1 close" data-dismiss="toast" aria-label="Close">
                  <span aria-hidden="true">&times;</span>
                </button>
              </div>
              <div class="toast-body">
                ${content}
              </div>
            </div>
          `;

    $('.toast').toast('show');
}

function appendTd(tr, content) {
    const td = document.createElement('td');
    td.appendChild(document.createTextNode(content));
    td.classList.add('white');
    tr.appendChild(td);
    return td
}

function appendButton(tr, content, onclick, ...classList) {
    const tdRemove = document.createElement('td');
    const btn = document.createElement('button');

    btn.type = 'submit';
    btn.classList.add('btn', 'btn-primary', 'btn-fill', 'mx-auto', ...classList);
    btn.appendChild(document.createTextNode(content));
    btn.onclick = onclick;

    tdRemove.appendChild(btn);
    tr.appendChild(tdRemove);
}

function appendButtonHref(tr, content, href) {
    const tdRemove = document.createElement('td');
    const btn = document.createElement('a');

    btn.href = href;
    btn.type = 'submit';
    btn.classList.add('btn', 'btn-primary', 'btn-fill', 'mx-auto');
    btn.appendChild(document.createTextNode(content));

    tdRemove.appendChild(btn);
    tr.appendChild(tdRemove);
}

function prependChild(parent, child) {
    if (parent.children.length === 0) {
        parent.appendChild(child);
    } else {
        parent.insertBefore(child, parent.children[0]);
    }
}

function createElement(tag, ...classList) {
    const el = document.createElement(tag);
    el.classList.add(...classList);
    return el;
}