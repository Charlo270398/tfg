function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    addLinkBreadcrumb('Historia clínica', '/user/patient/historial');
    loadTable(historial.entradas);
    cargarTablaAnaliticas(historial.analiticas);
}

function loadTable(eList){
    if(!eList || eList.length < 1){
        document.querySelector("#alert").textContent = "No hay ninguna entrada en tu historial";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector(`#historialTabla`).classList.add('invisible');
    }else{
        eList.forEach(entrada => {
            addRow(entrada);
        });
    }
}

function cargarTablaAnaliticas(aList){
    if(!aList || aList.length < 1){
        document.querySelector("#alertTablaAnaliticas").textContent = "No hay ninguna analítica en tu historial";
        document.querySelector("#alertTablaAnaliticas").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alertTablaAnaliticas").classList.remove('invisible');
        document.querySelector(`#analiticasTabla`).classList.add('invisible');
    }else{
        document.querySelector("#alert").classList.add('invisible');
        document.querySelector("#analiticasTabla").classList.remove('invisible');
        aList.forEach(analitica => {
            addRowAnaliticas(analitica);
        });
    }
}

function addRow(entrada){
    let tr = document.createElement('tr');
    let id = document.createElement('td');
    let fecha = document.createElement('td');
    let especialista = document.createElement('td');
    let tipo = document.createElement('td');
    let acciones = document.createElement('td');

    let accionesButton = document.createElement('button');
    accionesButton.classList = "btn btn-primary";
    accionesButton.type = "button";
    accionesButton.textContent = "Consultar entrada";
    accionesButton.addEventListener("click", consultarEntradaHistorial, false);
    acciones.append(accionesButton);
    fecha.textContent = entrada.createdAt;
    especialista.textContent = entrada.empleadoNombre;
    tipo.textContent = entrada.tipo;
    id.textContent = entrada.id;
    tr.append(id);
    tr.append(tipo);
    tr.append(especialista);
    tr.append(fecha);
    tr.append(acciones);
    tr.setAttribute("id", entrada.id);
    //Añadimos fila a la tabla
    document.querySelector(`#historialTabla`).querySelector('tbody').append(tr);
}

function addRowAnaliticas(entrada){
    let tr = document.createElement('tr');
    let id = document.createElement('td');
    let fecha = document.createElement('td');
    let especialista = document.createElement('td');
    let tipo = document.createElement('td');
    let acciones = document.createElement('td');

    let accionesButton = document.createElement('button');
    accionesButton.classList = "btn btn-primary";
    accionesButton.type = "button";
    accionesButton.textContent = "Consultar analítica";
    accionesButton.addEventListener("click", consultarAnaliticaHistorial, false);
    acciones.append(accionesButton);
    fecha.textContent = entrada.createdAt;
    especialista.textContent = entrada.empleadoNombre;
    tipo.textContent = "Analítica";
    id.textContent = entrada.id;
    tr.append(id);
    tr.append(tipo);
    tr.append(especialista);
    tr.append(fecha);
    tr.append(acciones);
    tr.setAttribute("id", entrada.id);
    //Añadimos fila a la tabla
    document.querySelector(`#analiticasTabla`).querySelector('tbody').append(tr);
}

function consultarEntradaHistorial(event){
    window.location.href = "/user/patient/historial/entrada?entradaId=" + event.target.closest("tr").getAttribute("id");
}

function consultarAnaliticaHistorial(event){
    window.location.href = "/user/patient/historial/analitica?analiticaId=" + event.target.closest("tr").getAttribute("id");
}

document.addEventListener('DOMContentLoaded',init,false);