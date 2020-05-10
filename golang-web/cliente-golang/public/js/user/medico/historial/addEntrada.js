function submit(event){

}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Medico', '/user/doctor');
    addLinkBreadcrumb('Solicitar historial', '/user/doctor/historial/solicitar');
    addLinkBreadcrumb('Añadir entrada', '');
    document.querySelector("#searchButton").addEventListener('click',submit,false);
}

document.addEventListener('DOMContentLoaded',init,false);