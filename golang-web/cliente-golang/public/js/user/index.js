function cargarTarjetasPermitidas(userRoles){
    userRoles.forEach(rol => {
        añadirTarjeta(rol);
    });
}
function añadirTarjeta(rol_id){
    document.querySelector("#rol_" + rol_id).classList.remove('invisible');
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    if(userRoles){
        cargarTarjetasPermitidas(userRoles);
    }else{
        console.log("ERROR CARGANDO COOKIE userRoles");
    }
}

document.addEventListener('DOMContentLoaded',init,false);