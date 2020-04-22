function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    addLinkBreadcrumb('Citas', '/user/patient/citas');
}

document.addEventListener('DOMContentLoaded',init,false);