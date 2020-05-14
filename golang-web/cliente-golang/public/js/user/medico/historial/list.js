function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Medico', '/user/doctor');
    addLinkBreadcrumb('Historiales disponibles', '/user/doctor/historial/list');
    loadTable(historialesList);
}

function loadTable(hList){
    console.log(hList);
    if(!hList || hList.length < 1){
        document.querySelector("#alert").textContent = "No hay ningun historial disponible";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
    hList.forEach(historial => {
        addRow(historial);
    });
}

function addRow(historial){
    let tr = document.createElement('tr');
    let nombre = document.createElement('td');
    let sexo = document.createElement('td');
    let consulta = document.createElement('td');

    let consultaButton = document.createElement('button');
    consultaButton.classList = "btn btn-primary";
    consultaButton.type = "button";
    consultaButton.textContent = "Consultar historial";
    consultaButton.addEventListener("click", consultarHistorial, false);
    consulta.append(consultaButton);
    nombre.textContent = historial.nombrePaciente + " " + historial.apellidosPaciente;
    sexo.textContent = historial.sexo;
    tr.append(nombre);
    tr.append(sexo);
    tr.append(consulta);
    tr.setAttribute("id", historial.id);
    //Añadimos fila a la tabla
    document.querySelector(`#historialesTabla`).querySelector('tbody').append(tr);
}
function consultarHistorial(event){
    window.location.href = "/user/doctor/historial?historialId=" + event.target.closest("tr").getAttribute("id");
}

document.addEventListener('DOMContentLoaded',init,false);