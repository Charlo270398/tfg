function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Emergencias', '/user/emergency');
    //Si se pasa por parametro el DNI se busca auto
    var url = new URL(window.location.href);
    var paramIdentificacion = url.searchParams.get("identificacion");
    if(paramIdentificacion){
        addLinkBreadcrumb('Historial', '/user/emergency?identificacion='+paramIdentificacion);
    }
    addLinkBreadcrumb('Consultar analítica', '');
}

document.addEventListener('DOMContentLoaded',init,false);