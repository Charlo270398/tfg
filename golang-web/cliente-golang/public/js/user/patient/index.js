function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    if(permisos){
        document.querySelector("#alert").classList.remove('invisible');
    }
}

document.addEventListener('DOMContentLoaded',init,false);