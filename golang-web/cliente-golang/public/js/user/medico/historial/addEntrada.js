var HISTORIAL_ID = -1;

function submit(event){
    if(document.querySelector("#motivoConsulta").value != "" && document.querySelector("#juicioDiagnostico").value != ""){
        console.log("AA");
        //restAddEntrada(document.querySelector("#motivoConsulta").value, document.querySelector("#juicioDiagnostico").value);
    }else{
        document.querySelector("#alert").textContent = "Existen campos vacíos";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
}

function restAddEntrada(motivoConsulta, juicioDiagnostico){
    const url= `/user/doctor/citas/addEntrada`;
    const payload= {citaId: parseInt(cita.id), pacienteId: parseInt(cita.pacienteId), motivoConsulta: motivoConsulta, juicioDiagnostico: juicioDiagnostico};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                //Cerrar cita
                console.log(r);
                if(r.Result == "OK"){
                    document.querySelector("#alert").textContent = "Entrada insertada correctamente";
                    document.querySelector("#alert").classList.remove('invisible');
                    document.querySelector("#formConsulta").classList.add('invisible');
                }
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
            }
        })
        .catch(err => alert(err));
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Médico', '/user/doctor');
    addLinkBreadcrumb('Solicitar historial', '/user/doctor/historial/solicitar');
    //Si se pasa por parametro el DNI se busca auto
    var url = new URL(window.location.href);
    var paramIdentificacion = url.searchParams.get("identificacion");
    if(paramIdentificacion){
        addLinkBreadcrumb('Historial', '/user/doctor/historial/solicitar?identificacion='+paramIdentificacion);
    }
    addLinkBreadcrumb('Añadir entrada', '');
    var paramHistorialId = url.searchParams.get("historialId");
    if(paramHistorialId){
        HISTORIAL_ID = parseInt(paramHistorialId);
    }
}

document.addEventListener('DOMContentLoaded',init,false);