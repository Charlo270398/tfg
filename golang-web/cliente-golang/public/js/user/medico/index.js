function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Medico', '/user/doctor');
    if(citaActualId != -1){
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector("#referenciaCita").setAttribute("href", "/user/doctor/citas?citaId=" + citaActualId);
    }
}

document.addEventListener('DOMContentLoaded',init,false);