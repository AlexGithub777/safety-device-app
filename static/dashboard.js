// hot reload
if (window.EventSource) {
    new EventSource("http://localhost:8090/internal/reload").onmessage = () => {
        location.reload();
    };
}

// Leaflet map setup
// SVG's intrinsic width/height or viewBox values
const svgWidth = 561.568;
const svgHeight = 962.941;
const minX = 128.009;
const minY = 82.331;

var map = L.map("map", {
    crs: L.CRS.Simple,
    minZoom: -1,
});

var imageUrl = "/static/map.svg";

// Apply these values to the bounds
const bounds = [
    [0, 0],
    [svgHeight, svgWidth],
];

L.imageOverlay(imageUrl, bounds).addTo(map);

// Fetch the buildings data draw overlay nat each co ordinate
fetch("static/buildings.json")
    .then((response) => response.json())
    .then((data) => {
        console.log(data);

        // Loop over the buildings and add a rectangle for each one
        data.buildings.forEach((building) => {
            var x = building.coordinates.x - minX; // Subtract the min-x value
            var y = svgHeight - (building.coordinates.y - minY); // Subtract the min-y value and flip the y-coordinate

            // Create the rectangle
            var rectangle = L.rectangle([
                [y - 19, x], // Subtract the height of the rectangle from the y-coordinate
                [y, x + 19],
            ]).addTo(map);

            // Add a click event listener to the rectangle
            rectangle.on("click", () => {
                // Fetch devices for the clicked building
                GetAllDevices(building.name); // Pass the building code to fetchDevices
                console.log("Building clicked:", building.name);
            });
        });
    });

map.fitBounds(bounds);

async function GetAllDevices(buildingCode = "") {
    try {
        const url = buildingCode
            ? `/api/emergency-device?building_code=${buildingCode}`
            : "/api/emergency-device";
        const response = await fetch(url);
        const devices = await response.json();

        const tbody = document.getElementById("emergency-device-body");
        tbody.innerHTML = devices
            .map((device) => {
                // Helper function to format dates as "Month YYYY" or "N/A" for manufacture and expiry dates
                const formatDateMonthYear = (dateString) => {
                    if (!dateString || dateString === "0001-01-01T00:00:00Z") {
                        return "N/A";
                    }
                    const date = new Date(dateString);
                    return date.toLocaleDateString("en-US", {
                        year: "numeric",
                        month: "long",
                    });
                };

                // Helper function to format dates as "Month Day, YYYY" or "N/A" for inspection dates
                const formatDateFull = (dateString) => {
                    if (!dateString || dateString === "0001-01-01T00:00:00Z") {
                        return "N/A";
                    }
                    const date = new Date(dateString);
                    return date.toLocaleDateString("en-US", {
                        year: "numeric",
                        month: "long",
                        day: "numeric",
                    });
                };

                // Get the badge class based on the status
                const badgeClass =
                    device.status.String === "Active"
                        ? "text-bg-success"
                        : device.status.String === "Expired"
                        ? "text-bg-danger"
                        : "text-bg-warning";

                // Return the row with formatted dates
                // Return the row with formatted dates and data-* attributes
                return `
                    <tr>
                        <td data-label="Device Type">${
                            device.emergency_device_type_name
                        }</td>
                        <td data-label="Extinguisher Type">${
                            device.extinguisher_type_name.String
                        }</td>
                        <td data-label="Room">${device.room_code}</td>
                        <td data-label="Serial Number">${
                            device.serial_number.String
                        }</td>
                        <td data-label="Manufacture Date">${formatDateMonthYear(
                            device.manufacture_date.Time
                        )}</td>
                        <td data-label="Expire Date">${formatDateMonthYear(
                            device.expire_date.Time
                        )}</td>
                        <td data-label="Last Inspection Date">${formatDateFull(
                            device.last_inspection_date.Time
                        )}</td>
                        <td data-label="Next Inspection Date">${formatDateFull(
                            device.next_inspection_date.Time
                        )}</td>
                        <td data-label="Size">${device.size.String}</td>
                        <td data-label="Status"><span class="badge ${badgeClass}">${
                    device.status.String
                }</span></td>
            <td>
                <div class="action-buttons">
                                <button class="btn btn-secondary" onclick="ViewDeviceInspection(${
                                    device.emergency_device_id
                                })">Inspect</button>
                                <button class="btn btn-primary" onclick="DeviceNotes('${
                                    device.description.String
                                }')">Notes</button>
                                <button class="btn btn-warning" onclick="EditDevice(${
                                    device.emergency_device_id
                                })">Edit</button>
                                <button class="btn btn-danger" onclick="DeleteDevice(${
                                    device.emergency_device_id
                                })">Delete</button>
                            </div>
                        </td>
                    </tr>
                `;
            })
            .join("");
    } catch (err) {
        console.error("Failed to fetch devices:", err);
    }
}

// Initial fetch without filtering
GetAllDevices();

function AddDevice() {
    // Fetch the device types and populate the select options
    fetch("/api/emergency-device-type")
        .then((response) => response.json())
        .then((data) => {
            const select = document.getElementById("status");
            // Clear previous options
            select.innerHTML = "";
            // Add a default option and select it
            const defaultOption = document.createElement("option");
            defaultOption.text = "Select Device type";
            defaultOption.selected = true;
            defaultOption.disabled = true;
            select.add(defaultOption);
            data.forEach((item) => {
                const option = document.createElement("option");
                option.text = item.emergency_device_type_name; // Set the text of the option
                select.add(option);
            });
        })
        .catch((error) => console.error("Error:", error));

    // Fetch the extinguisher types and populate the select options
    fetch("/api/extinguisher-type")
        .then((response) => response.json())
        .then((data) => {
            const select = document.getElementById("extinguisherType");
            // Clear previous options
            select.innerHTML = "";
            // Add a default option and select it
            const defaultOption = document.createElement("option");
            defaultOption.text = "Select Extinguisher Type";
            defaultOption.selected = true;
            defaultOption.disabled = true;
            select.add(defaultOption);
            data.forEach((item) => {
                const option = document.createElement("option");
                option.text = item.extinguisher_type_name; // Set the text of the option
                select.add(option);
            });
        })
        .catch((error) => console.error("Error:", error));

    // Function to populate the room dropdown
    function populateRooms(buildingId) {
        // Fetch the rooms for the selected building and populate the select options
        fetch(`/api/room?buildingId=${buildingId}`)
            .then((response) => response.json())
            .then((data) => {
                const select = document.getElementById("room");
                // Clear previous options
                select.innerHTML = "";
                // Add a default option and select it
                const defaultOption = document.createElement("option");
                defaultOption.text = "Select Room";
                defaultOption.selected = true;
                defaultOption.disabled = true;
                select.add(defaultOption);
                data.forEach((item) => {
                    const option = document.createElement("option");
                    option.text = item.room_code; // Set the text of the option
                    select.add(option);
                });
            })
            .catch((error) => console.error("Error:", error));
    }

    // Fetch the buildings and populate the select options
    fetch("/api/building")
        .then((response) => response.json())
        .then((data) => {
            const select = document.getElementById("building");
            // Clear previous options
            select.innerHTML = "";
            // Add a default option and select it
            const defaultOption = document.createElement("option");
            defaultOption.text = "Select Building";
            defaultOption.selected = true;
            defaultOption.disabled = true;
            select.add(defaultOption);
            data.forEach((item) => {
                const option = document.createElement("option");
                option.text = item.building_code; // Set the text of the option
                select.add(option);
            });
            // Populate the room dropdown for the first building
            if (data.length > 0) {
                populateRooms(data[0].building_code);
            }
            // Add event listener to the building select
            select.addEventListener("change", function () {
                populateRooms(this.value);
            });
        })
        .catch((error) => console.error("Error:", error));

    // Show the modal after populating the select options
    $("#addModal").modal("show");
}

function EditDevice(deviceId) {
    console.log(`Edit device with ID: ${deviceId}`);
    // Add your edit logic here
}

function DeleteDevice(deviceId) {
    console.log(`Delete device with ID: ${deviceId}`);
    // Add your delete logic here
}

// Change to add inspection
function ViewDeviceInspection(deviceId) {
    console.log(`Inspect device with ID: ${deviceId}`);

    // Show the modal
    $("#viewInspectionModal").modal("show");
}

function ViewInspectionDetails(inspectionId) {
    console.log(`View inspection details for inspection ID: ${inspectionId}`);
    // Add your view inspection details logic here
}

function AddInspection() {
    // Close the view inspection modal
    $("#viewInspectionModal").modal("hide");

    // Show the modal
    $("#addInspectionModal").modal("show");
}

function DeviceNotes(description) {
    // Populate the modal with the description
    document.getElementById("notesModalBody").innerText = description;

    // Show the modal
    $("#notesModal").modal("show");
}

function ViewNotifications() {
    console.log("View notifications");
    $("#notificationsModal").modal("show");
    // Add your view notifications logic here
}

// Function to toggle the map visibility
function ToggleMap() {
    var map = document.getElementById("map");
    var deviceList = document.querySelector(".device-list");

    // Check if the map is currently visible
    if (map.classList.contains("d-none")) {
        // Map is hidden, show the map and set device list back to col-xxl-9 width
        map.classList.remove("d-none");
        map.classList.add("col-xxl-2");
        deviceList.classList.remove("col-xxl-12");
        deviceList.classList.add("col-xxl-10");
    } else {
        // Map is visible, hide the map and make device list 100% width
        map.classList.add("d-none");
        deviceList.classList.remove("col-xxl-10");
        deviceList.classList.add("col-xxl-12");
    }
}
