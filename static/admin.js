// hot reload
if (window.EventSource) {
    new EventSource("http://localhost:8090/internal/reload").onmessage = () => {
        setTimeout(() => {
            location.reload();
        });
    };
}
// Fetch users from the server
// Fetch site data from the server
fetch("/api/site")
    .then((response) => response.json())
    .then((sites) => {
        // Create a table row for each site
        const siteRows = sites.map(
            (site) => `
<tr>
<td>${site.site_name}</td>
<td>${site.site_address}</td>
<td>
    <div class="btn-group">
        <button class="btn btn-primary edit-button" data-id="${site.site_id}">Edit</button>
        <button class="btn btn-danger delete-button" data-id="${site.site_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the sites table
        $("#sites-table tbody").html(siteRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle edit
        });
        $(".delete-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle delete
        });
    });

// Fetch buildings from the server
fetch("/api/building")
    .then((response) => response.json())
    .then((buildings) => {
        // Create a table row for each building
        const buildingRows = buildings.map(
            (building) => `
<tr>
<td>${building.building_code}</td>
<td>${building.site_name}</td>
<td>
    <div class="btn-group">
        <button class="btn btn-primary edit-button" data-id="${building.building_id}">Edit</button>
        <button class="btn btn-danger delete-button" data-id="${building.building_id}">Delete</button>
    </div> 
</td>
</tr>
`
        );

        // Add the rows to the buildings table
        $("#buildings-table tbody").html(buildingRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle edit
        });
        $(".delete-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle delete
        });
    });

// Fetch rooms from the server
fetch("/api/room")
    .then((response) => response.json())
    .then((rooms) => {
        // Create a table row for each room
        const roomRows = rooms.map(
            (room) => `
<tr>
<td>${room.room_code}</td>
<td>${room.building_code}</td>
<td>${room.site_name}</td>
<td>
    <div class="btn-group">
        <button class="btn btn-primary edit-button" data-id="${room.room_id}">Edit</button>
        <button class="btn btn-danger delete-button" data-id="${room.room_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the rooms table
        $("#rooms-table tbody").html(roomRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle edit
        });
        $(".delete-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle delete
        });
    });

// Fetch device types from the server
fetch("/api/emergency-device-type")
    .then((response) => response.json())
    .then((deviceTypes) => {
        // Create a table row for each device type
        const deviceTypeRows = deviceTypes.map(
            (deviceType) => `
<tr>
<td>${deviceType.emergency_device_type_name}</td>
<td>
    <div class="btn-group">
        <button class="btn btn-primary edit-button" data-id="${deviceType.emergency_device_type_id}">Edit</button>
        <button class="btn btn-danger delete-button" data-id="${deviceType.emergency_device_type_id}">Delete</button>
    </div>
</td>
</tr>
`
        );

        // Add the rows to the device types table
        $("#device-types-table tbody").html(deviceTypeRows.join(""));

        // Add event listeners to the edit and delete buttons
        $(".edit-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle edit
        });
        $(".delete-button").click((event) => {
            const id = $(event.target).data("id");
            // Handle delete
        });
    });