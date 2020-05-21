var HISTORIAL_ID;

function submit(event){
    if(document.querySelector("#motivoConsulta").value != "" && document.querySelector("#juicioDiagnostico").value != "" && 
    document.querySelector("#tipoSelector").value != "-1" && HISTORIAL_ID != undefined){
        restAddEntrada(document.querySelector("#motivoConsulta").value, document.querySelector("#juicioDiagnostico").value, document.querySelector("#tipoSelector").value);
    }else{
        document.querySelector("#alert").textContent = "Existen campos vacíos";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
}

function restAddEntrada(motivoConsulta, juicioDiagnostico, tipo){
    const url= `/user/emergency/historial/addEntrada`;
    const payload= {historialId: HISTORIAL_ID, motivoConsulta: motivoConsulta, juicioDiagnostico: juicioDiagnostico, tipo: tipo};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
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
    addLinkBreadcrumb('Emergencias', '/user/emergency');
    //Si se pasa por parametro el DNI se busca auto
    var url = new URL(window.location.href);
    var paramIdentificacion = url.searchParams.get("identificacion");
    var paramHistorialId = url.searchParams.get("historialId");
    if(paramIdentificacion){
        addLinkBreadcrumb('Historial', '/user/emergency?identificacion='+paramIdentificacion);
    }
    if(paramHistorialId){
        HISTORIAL_ID = parseInt(paramHistorialId);
    }
    addLinkBreadcrumb('Añadir entrada', '');
    document.querySelector("#searchButton").addEventListener('click',submit,false);
}

document.addEventListener('DOMContentLoaded',init,false);