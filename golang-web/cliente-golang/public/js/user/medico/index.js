function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Medico', '/user/doctor');
    if(citaActualId != -1){
        document.querySelector("#alert").classList.remove('invisible');
    }
}

document.addEventListener('DOMContentLoaded',init,false);