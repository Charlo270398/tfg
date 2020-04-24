function submit(event){
    let nombre = document.querySelector("#name").value;
    let apellido1 = document.querySelector("#surname1").value;
    let apellido2 = document.querySelector("#surname2").value;
    let identificacion = document.querySelector("#idnumber").value;
    let email = document.querySelector("#email").value;
    let password = document.querySelector("#pass").value;
    let especialidadMedico = document.querySelector("#medicoEspecialidadSelector").value;
    if(nombre && apellido1 && apellido2 && identificacion && email && password){
        if(clinicaId == -1 || especialidadMedico == -1){
            document.querySelector("#alert").textContent = "Debes rellenar los datos para ejercer como médico";
            document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
            document.querySelector("#alert").classList.remove('invisible');
            return;
        }
        register(nombre,apellido1, apellido2, email, identificacion, password, clinicaId.toString(), especialidadMedico);
    }else{
        document.querySelector("#alert").textContent = "Existen campos vacíos";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
}
function register(nombre, apellido1, apellido2, email, identificacion, password, medicoClinica, medicoEspecialidad){
    var result = false;
    let apellidos = apellido1;
    if(apellido2){
        apellidos += " " + apellido2;
    }
    const url= `/user/admin/doctor/add`;
    const payload= {nombre: nombre, identificacion:identificacion, apellidos: apellidos, email: email, password: password,
         medicoClinica: medicoClinica, medicoEspecialidad: medicoEspecialidad};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                console.log("USUARIO AÑADIDO CORRECTAMENTE");
                document.querySelector("#alert").classList.replace("alert-danger", "alert-success");
                document.querySelector("#alert").textContent = "Usuario añadido correctamente";
                document.querySelector("#alert").classList.remove('invisible');
                //window.location.href="/user/menu";
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
            }
        })
        .catch(err => alert(err));
    return result;
}

function loadEspecialidades(especialidades){
    var medicoClinicaSelector = document.querySelector("#medicoEspecialidadSelector");
    especialidades.forEach(e => {
        var optionE = document.createElement("option");
        optionE.value = e.id;
        optionE.textContent = e.nombre;
        medicoClinicaSelector.append(optionE);
    });
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Administrador', '/user/admin');
    addLinkBreadcrumb('Añadir medico', '/user/admin/doctor/add');
    document.querySelector("#submit").addEventListener('click',submit,false);
    loadEspecialidades(especialidades);
}

document.addEventListener('DOMContentLoaded',init,false);