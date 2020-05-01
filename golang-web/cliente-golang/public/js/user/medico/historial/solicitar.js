function busquedaDNI(event){
    if(document.querySelector("#inputDNI").value.length != 9){
        //Activar alerta
        document.querySelector("#alert").textContent = "El documento de identificación debe tener un formato válido (por ejemplo, 00000000X)";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector("#historialTabla").classList.add('invisible');
        document.querySelector("#historialTabla").classList.add('invisible');
        document.querySelector("#historialTabla").classList.add('buttonsForm');
        return;
    }else{
        restBuscarDNI(document.querySelector("#inputDNI").value);
    }
}

function cargarTablaHistorial(historial){
    document.querySelector("#alert").classList.add('invisible');
    document.querySelector("#historialTabla").classList.remove('invisible');
}

function restBuscarDNI(DNI){
    const url= `/user/doctor/historial/solicitar`;
    const payload= {identificacion: DNI};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                //PROCESAR HISTORIAL
                cargarTablaHistorial(r);
                document.querySelector("#buttonsForm").classList.remove('invisible');
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
                document.querySelector("#historialTabla").classList.add('invisible');
            }
        })
        .catch(err => alert(err));
}

function addEntrada(){
    window.location.href = "/user/doctor/historial/addEntrada";
}

function addAnalitica(){
    window.location.href = "/user/doctor/historial/addAnalitica";
}

function solicitarAccesoTotal(){

}

function solicitarAccesoEntrada(){

}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Medico', '/user/doctor');
    addLinkBreadcrumb('Solicitar historial', '/user/historial/solicitar');
    document.querySelector("#searchButton").addEventListener('click',busquedaDNI,false);
    document.querySelector("#accesoTotalButton").addEventListener('click',solicitarAccesoTotal,false);
    document.querySelector("#addEntradaButton").addEventListener('click',addEntrada,false);
    document.querySelector("#addAnaliticaButton").addEventListener('click',addAnalitica,false);
}

document.addEventListener('DOMContentLoaded',init,false);