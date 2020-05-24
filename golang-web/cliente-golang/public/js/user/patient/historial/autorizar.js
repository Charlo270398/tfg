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
            autorizarButton.addEventListener("click", autorizarSolicitudEntrada, false);
            denegarButton.addEventListener("click", denegarSolicitudEntrada, false);
        }
        else if(solicitud.analiticaId != 0){
            tr.setAttribute("id", solicitud.analiticaId);
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
    console.log(event.target.closest("tr").getAttribute("tipo"));
}

function autorizarSolicitudEntrada(event){
    console.log(event.target.closest("tr").getAttribute("id"));
}

function autorizarSolicitudAnalitica(event){
    console.log(event.target.closest("tr").getAttribute("id"));
}

function denegarSolicitudHistorial(event){
    console.log(event.target.closest("tr").getAttribute("tipo"));
}

function denegarSolicitudEntrada(event){
    console.log(event.target.closest("tr").getAttribute("id"));
}

function denegarSolicitudAnalitica(event){
    console.log(event.target.closest("tr").getAttribute("id"));
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