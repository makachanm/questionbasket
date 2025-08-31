

window.onload = () => {
    fetch("http://localhost:3000/api/profile", {
        method: "GET"
    }).then(response => response.json()).then(data => {
        const title_name_holder = document.getElementsByClassName("profile-name");
        const title_description_holder = document.getElementsByClassName("profile-description");

        const name = data["name"];
        const description = data["description"];

        title_name_holder.item(0).innerHTML = name;
        title_description_holder.item(0).innerHTML = description;
    });
}