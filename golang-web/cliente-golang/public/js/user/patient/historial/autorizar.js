function loadTable(sList){
    if(!sList || sList.length < 1){
        document.querySelector("#alert").classList.remove('invisible');
    }else{
        sList.forEach(solicitud => {
            addRow(solicitud);
        });
    }
}

function addRow(solicitud){
    let tr = document.createElement('tr');
    let solicitante = document.createElement('td');
    let tipo = document.createElement('td');
    let idEntrada = document.createElement('td');
    let acciones = document.createElement('td');

    let autorizarButton = document.createElement('button');
    autorizarButton.classList = "btn btn-primary";
    autorizarButton.type = "button";
    autorizarButton.textContent = "Autorizar";

    let denegarButton = document.createElement('button');
    denegarButton.classList = "btn btn-danger";
    denegarButton.type = "button";
    denegarButton.textContent = "Denegar";

    solicitante.textContent = solicitud.nombreEmpleado;
    tr.setAttribute("empleadoId", solicitud.empleadoId);
    if(solicitud.tipoHistorial == "TOTAL"){
        tipo.textContent = "Acceso total al historial";
        idEntrada.textContent = "";
        tr.setAttribute("tipo", "TOTAL");
        autorizarButton.addEventListener("click", autorizarSolicitudHistorial, false);
        denegarButton.addEventListener("click", denegarSolicitudHistorial, false);
    }
    else if(solicitud.tipoHistorial == "BASICO"){
        tipo.textContent = "Acceso básico al historial";
        idEntrada.textContent = "";
        tr.setAttribute("tipo", "BASICO");
        autorizarButton.addEventListener("click", autorizarSolicitudHistorial, false);
        denegarButton.addEventListener("click", denegarSolicitudHistorial, false);
    }
    else{
        if(solicitud.entradaId != 0){
            tr.setAttribute("id", solicitud.entradaId);
            tipo.textContent = "Acceso a entrada";
            let link = document.createElement('a');
            link.href = "/user/patient/historial/entrada?entradaId=" +  solicitud.entradaId;
            link.textContent = "Entrada con ID " + solicitud.entradaId;
            idEntrada.append(link);
            autorizarButton.addEventListener("click", autorizarSolicitudEntrada, false);
            denegarButton.addEventListener("click", denegarSolicitudEntrada, false);
        }
        else if(solicitud.analiticaId != 0){
            tr.setAttribute("id", solicitud.analiticaId);
            tipo.textContent = "Acceso a analítica";
            let link = document.createElement('a');
            link.href = "/user/patient/historial/analitica?analiticaId=" +  solicitud.analiticaId;
            link.textContent = "Analítica con ID " + solicitud.analiticaId;
            idEntrada.append(link);
            autorizarButton.addEventListener("click", autorizarSolicitudAnalitica, false);
            denegarButton.addEventListener("click", denegarSolicitudAnalitica, false);
        }
    }
    acciones.append(autorizarButton);
    acciones.append(denegarButton);
    tr.append(solicitante);
    tr.append(tipo);
    tr.append(idEntrada);
    tr.append(acciones);
    //Añadimos fila a la tabla
    document.querySelector(`#solicitudesTablaBody`).append(tr);
}

function autorizarSolicitudHistorial(event){
    permitir(event.target.closest("tr").getAttribute("empleadoId"),event.target.closest("tr").getAttribute("tipo"), null, null, event.target);
}

function autorizarSolicitudEntrada(event){
    permitir(event.target.closest("tr").getAttribute("empleadoId"), null, event.target.closest("tr").getAttribute("id"), null, event.target);
}

function autorizarSolicitudAnalitica(event){
    permitir(event.target.closest("tr").getAttribute("empleadoId"), null, null, event.target.closest("tr").getAttribute("id"), event.target);
}

function denegarSolicitudHistorial(event){
    denegar(event.target.closest("tr").getAttribute("empleadoId"),event.target.closest("tr").getAttribute("tipo"), null, null, event.target);
}

function denegarSolicitudEntrada(event){
    denegar(event.target.closest("tr").getAttribute("empleadoId"),null, event.target.closest("tr").getAttribute("id"), null, event.target);
}

function denegarSolicitudAnalitica(event){
    denegar(event.target.closest("tr").getAttribute("empleadoId"),null, null, event.target.closest("tr").getAttribute("id"), event.target);
}

function denegar(empleadoId, tipo, entradaId, analiticaId, button){
    //Enviamos peticion
    const url= `/permisos/solicitudes/denegar`;
    const payload= {empleadoId: parseInt(empleadoId), tipoHistorial: tipo, entradaId: parseInt(entradaId), analiticaId: parseInt(analiticaId)};
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
                    acciones = button.closest("td");
                    console.log("SOLICITUD BORRADA");
                    //Eliminamos datos anteriores
                    while (acciones.lastElementChild) { 
                        acciones.removeChild(acciones.lastElementChild); 
                    } 
                    let text = document.createElement('span');
                    text.textContent = "SOLICITUD DENEGADA";
                    acciones.append(text);
                }
            }
        })
    .catch(err => alert(err));
}

function permitir(empleadoId, tipo, entradaId, analiticaId, button){
    //Enviamos peticion
    const url= `/permisos/solicitudes/permitir`;
    const payload= {empleadoId: parseInt(empleadoId), tipoHistorial: tipo, entradaId: parseInt(entradaId), analiticaId: parseInt(analiticaId)};
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
                    acciones = button.closest("td");
                    console.log("SOLICITUD ACEPTADA");
                    //Eliminamos datos anteriores
                    while (acciones.lastElementChild) { 
                        acciones.removeChild(acciones.lastElementChild); 
                    } 
                    let text = document.createElement('span');
                    text.textContent = "SOLICITUD ACEPTADA";
                    acciones.append(text);
                }
            }
        })
    .catch(err => alert(err));
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    addLinkBreadcrumb('Autorizar', '');
    if(solicitudes){
        loadTable(solicitudes);
    }else{
        document.querySelector("#alert").classList.remove('invisible');
    }
}

document.addEventListener('DOMContentLoaded',init,false);