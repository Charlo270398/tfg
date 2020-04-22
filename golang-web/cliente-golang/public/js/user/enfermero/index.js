function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Enfermero', '/user/nurse');
}

document.addEventListener('DOMContentLoaded',init,false);