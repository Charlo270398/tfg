function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    addLinkBreadcrumb('Historia clínica', '/user/patient/historial');
}

document.addEventListener('DOMContentLoaded',init,false);