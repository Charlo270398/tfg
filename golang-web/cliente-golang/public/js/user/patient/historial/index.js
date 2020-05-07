function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    addLinkBreadcrumb('Historia clínica', '/user/patient/historial');
    console.log(historial);
    loadTable(historial.entradas)
}

function loadTable(eList){
    if(!eList || eList.length < 1){
        document.querySelector("#alert").textContent = "No hay ninguna entrada en tu historial";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector(`#historialTabla`).classList.add('invisible');
        document.querySelector(`#historialTitulo`).classList.add('invisible');
    }else{
        eList.forEach(entrada => {
            addRow(entrada);
        });
    }
}

function addRow(entrada){
    let tr = document.createElement('tr');
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
    tr.append(tipo);
    tr.append(especialista);
    tr.append(fecha);
    tr.append(acciones);
    tr.setAttribute("id", historial.id);
    //Añadimos fila a la tabla
    document.querySelector(`#historialTabla`).querySelector('tbody').append(tr);
}

function consultarEntradaHistorial(event){
    window.location.href = "/user/patient/historial/entrada?entradaId=" + event.target.closest("tr").getAttribute("id");
}

document.addEventListener('DOMContentLoaded',init,false);