function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Administrador clínica', '/user/admin');
}

document.addEventListener('DOMContentLoaded',init,false);