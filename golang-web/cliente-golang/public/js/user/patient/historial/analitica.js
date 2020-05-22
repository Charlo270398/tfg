function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    addLinkBreadcrumb('Historial', '/user/patient/historial');
    addLinkBreadcrumb('Consultar analítica', '');
}

document.addEventListener('DOMContentLoaded',init,false);