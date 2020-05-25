function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Medico', '/user/doctor');
    //Si se pasa por parametro el DNI se busca auto
    var url = new URL(window.location.href);
    var paramIdentificacion = url.searchParams.get("identificacion");
    if(paramIdentificacion){
        addLinkBreadcrumb('Solicitar historial', '/user/doctor/historial/solicitar?identificacion='+paramIdentificacion);
    }
    addLinkBreadcrumb('Consultar entrada', '');
}

document.addEventListener('DOMContentLoaded',init,false);