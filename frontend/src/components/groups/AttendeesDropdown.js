import React, { useState } from "react";
import "../../styles/groups/EventModal.css";

export default function AttendeesDropdown({
  eventUUID,
  attendCount,
  attendees = [],
}) {
  const [showDropdown, setShowDropdown] = useState(false);

  const toggleDropdown = () => {
    setShowDropdown(!showDropdown);
  };

  const acceptedAttendees = attendees.filter((a) => a.response === "accepted");

  return (
    <div className="attendees-dropdown">
      <p>
        <strong>Attending:</strong>
        <button onClick={toggleDropdown} className="attendees-toggle-btn">
          {attendCount} going ▼
        </button>
      </p>

      {showDropdown && (
        <div className="attendees-dropdown-menu">
          {acceptedAttendees.length === 0 ? (
            <div className="attendees-empty-state">No attendees yet</div>
          ) : (
            acceptedAttendees.map((attendee, idx) => (
              <div key={idx} className="attendee-item">
                {attendee.created_by_name}
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}
