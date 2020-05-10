var enfermero = false;
var medico = false;
var admin = false;

function submit(event){
    let nombre = document.querySelector("#name").value;
    let apellido1 = document.querySelector("#surname1").value;
    let apellido2 = document.querySelector("#surname2").value;
    let identificacion = document.querySelector("#idnumber").value;
    let email = document.querySelector("#email").value;
    let password = document.querySelector("#pass").value;
    let clinicaEnfermero = document.querySelector("#enfermeroClinicaSelector").value;
    let clinicaMedico = document.querySelector("#medicoClinicaSelector").value;
    let especialidadMedico = document.querySelector("#medicoEspecialidadSelector").value;
    let clinicaAdmin = document.querySelector("#adminClinicaSelector").value;
    let roles = rolesSeleccionados();
    if(nombre && apellido1 && apellido2 && identificacion && email && password && roles.length > 0){
        if(enfermero == true && clinicaEnfermero == -1){
            document.querySelector("#alert").textContent = "Debes elegir una clínica donde ejercer como enfermero";
            document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
            document.querySelector("#alert").classList.remove('invisible');
            return;
        }
        if(medico == true && (clinicaMedico == -1 || especialidadMedico == -1)){
            document.querySelector("#alert").textContent = "Debes rellenar los datos para ejercer como médico";
            document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
            document.querySelector("#alert").classList.remove('invisible');
            return;
        }
        if(admin == true && clinicaAdmin == -1){
            document.querySelector("#alert").textContent = "Debes elegir una clínica a administrar";
            document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
            document.querySelector("#alert").classList.remove('invisible');
            return;
        }   
        register(nombre,apellido1, apellido2, email, identificacion, password, roles, clinicaEnfermero, clinicaMedico, clinicaAdmin, especialidadMedico);
    }else{
        document.querySelector("#alert").textContent = "Existen campos vacíos";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
}
function register(nombre, apellido1, apellido2, email, identificacion, password, roles, enfermeroClinica, medicoClinica, adminClinica, medicoEspecialidad){
    var result = false;
    let apellidos = apellido1;
    if(apellido2){
        apellidos += " " + apellido2;
    }
    const url= `/user/adminG/userList/add`;
    const payload= {nombre: nombre, identificacion:identificacion, apellidos: apellidos, email: email, password: password, roles: roles
    , enfermeroClinica: enfermeroClinica, medicoClinica: medicoClinica, adminClinica: adminClinica, medicoEspecialidad: medicoEspecialidad};
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
                document.querySelector("#submit").classList.add('invisible');
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

function rolesSeleccionados(){
    const rolesList = [];
    const divList = document.querySelector("checkbox-roles").shadowRoot.querySelector("#rolesGroup").children;
    for(i=0; i<divList.length; i++){
        if(divList[i].querySelector('input').checked){
            rolesList.push(parseInt(divList[i].getAttribute('rol_id')));
        }
    }
    return rolesList;
}

function mostrarEnfermeroClinicaGroup(event){
    document.querySelector("#enfermeroClinicaGroup").classList.toggle("invisible");
    enfermero = !enfermero;
}

function mostrarMedicoClinicaGroup(event){
    document.querySelector("#medicoClinicaGroup").classList.toggle("invisible");
    document.querySelector("#medicoEspecialidadGroup").classList.toggle("invisible");
    medico = !medico;
}

function mostrarAdminClinicaGroup(event){
    document.querySelector("#adminClinicaGroup").classList.toggle("invisible");
    admin = !admin;
}

function loadClinicas(clinicas){
    var enfermeroClinicaSelector = document.querySelector("#enfermeroClinicaSelector");
    var medicoClinicaSelector = document.querySelector("#medicoClinicaSelector");
    var adminClinicaSelector = document.querySelector("#adminClinicaSelector");
    clinicas.forEach(c => {
        var optionE = document.createElement("option");
        optionM = optionE.cloneNode()
        optionA = optionE.cloneNode()
        optionE.value = c.id;
        optionE.textContent = c.nombre;
        optionM.value = c.id;
        optionM.textContent = c.nombre;
        optionA.value = c.id;
        optionA.textContent = c.nombre;
        enfermeroClinicaSelector.append(optionE);
        medicoClinicaSelector.append(optionM);
        adminClinicaSelector.append(optionA);
    });
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
    addLinkBreadcrumb('Administrador global', '/user/adminG');
    addLinkBreadcrumb('Usuarios', '/user/adminG/userList');
    addLinkBreadcrumb('Añadir', '/user/adminG/userList/add');
    document.querySelector("#submit").addEventListener('click',submit,false);
    loadClinicas(clinicas);
    loadEspecialidades(especialidades);
}

document.addEventListener('DOMContentLoaded',init,false);