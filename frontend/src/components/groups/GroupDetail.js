// src/components/groups/GroupDetail.js
// groups/uuid/main-view, main-feed

import React, { useState, useEffect } from "react";
import { toast } from "react-toastify";

import Modal from "../Modal";
import CreatePostForm from "./[uuid]/posts/CreatePostForm";
import CreateEventForm from "./[uuid]/events/CreateEventForm";
import GroupHeader from "./[uuid]/GroupHeader";
import EventCard from "./[uuid]/events/EventCard";
import { createPost } from "../../lib/apiPosts";
import PostList from "../../components/PostList";
import AttendeesDropdown from "./AttendeesDropdown";

import { formatDate } from "../../utils/formatDate";

import "../../styles/groups/GroupDetail.css";
import "../../styles/groups/EventModal.css";

export default function GroupDetail(
  {   group, 
      members, 
      requests,
      allUsers,
      loadingUsers,
      onBack, 
      onRequestJoin,
      onApproveRequest,   
      onDenyRequest,
      onInviteUser,
      handleAcceptInvite,
      handleDeclineInvite,
      loadingActions = {},
      isJoining = false,
      isInviting = false,
      posts
  }) {

  if (!group) return null;

  // Modal state
  const [showPostModal, setShowPostModal] = useState(false);
  const [showEventModal, setShowEventModal] = useState(false);
  const [showMembersModal, setShowMembersModal] = useState(false);

  const handleOpenMembersModal = () => setShowMembersModal(true);
  const handleCloseMembersModal = () => setShowMembersModal(false);
  
  // Events state
  const [groupEvents, setGroupEvents] = useState([]);
  const [loadingEvents, setLoadingEvents] = useState(true);

  // Event modal state
  const [selectedEvent, setSelectedEvent] = useState(null);
  const [currentRSVP, setCurrentRSVP] = useState(null);
  const [eventAttendees, setEventAttendees] = useState([]);

  // Get current user UUID
  const currentUserUUID =
    typeof window !== "undefined" ? localStorage.getItem("user-uuid") : null;
  const isCreator = group.creator_uuid === currentUserUUID;

    // status checks
    const isMember = group.status === "accepted";
    const isRequested = group.status === "requested";
    const isInvited = group.status === "invited";

    // filter members and requests
    const acceptedMembers = members.filter(m => m.status === 'accepted');
    const pendingRequests = requests.filter(r => r.status === 'requested' || r.status === 'invited');


    // const memberUuids = new Set(members.map(m => m.follower_uuid || m.uuid));    
    // const nonMembers = allUsers.filter(u => !memberUuids.has(u.uuid));

    const memberUuids = new Set([
      ...members.map(m => m.follower_uuid || m.uuid),
      ...requests.map(r => r.follower_uuid || r.uuid)
    ]);
    const nonMembers = allUsers.filter(u => !memberUuids.has(u.uuid) && u.uuid !== currentUserUUID);
    const filteredNonMembers = nonMembers.filter(u => u.uuid !== currentUserUUID);

    const refreshEvents = () => {
        fetch(`/frontend-api/groups/events/${group.uuid}`)
            .then(res => {
            if (!res.ok) {
                return res.text().then(text => {
                    throw new Error(`HTTP error! Status: ${res.status}. Response: ${text}`);
                });
            }
            return res.json();
        })
        .then(data => setGroupEvents(data.data || []))
        .catch(err => {
            console.error("Failed to refresh events:", err);
            toast.error("Failed to load events");
        });
    };

    // Fetch events for this group
    useEffect(() => {
        if (!group?.uuid) return;
        setLoadingEvents(true);
        fetch(`/frontend-api/groups/events/${group.uuid}`)
            .then(res => res.json())
            .then(data => {
                setGroupEvents(data.data || []);
                setLoadingEvents(false);
            })
            .catch(() => setLoadingEvents(false));
    }, [group?.uuid]);

  // Check if current user is going to the event
  const isUserGoing = eventAttendees.some(
    (attendee) =>
      attendee.created_by_uuid === currentUserUUID &&
      attendee.response === "accepted"
  );

  // Fetch event attendees
  const fetchAttendees = async (eventUUID) => {
    try {
      const response = await fetch(
        `/frontend-api/group/event/responses/${eventUUID}`,
        {
          credentials: "include",
        }
      );
      if (response.ok) {
        const data = await response.json();
        setEventAttendees(data.data || []);
      }
    } catch (err) {
      console.error("Failed to fetch attendees:", err);
    }
  };

  // Fetch events for this group
  useEffect(() => {
    if (!group?.uuid) return;
    setLoadingEvents(true);
    fetch(`/frontend-api/groups/events/${group.uuid}`)
      .then((res) => res.json())
      .then((data) => {
        setGroupEvents(data.data || []);
        setLoadingEvents(false);
      })
      .catch(() => setLoadingEvents(false));
  }, [group?.uuid]);

  const handlePostSubmit = async (postData) => {
    try {
      const data = await createPost(postData);
      if (data) {
        window.location.href = `/post/${data.data}`;
      } else {
        toast.error(data.message || "Failed to create post");
      }
    } catch (err) {
      toast.error(err.message || "Error creating post");
    }
    toast.success(`Post created: ${postData.title}`);
    setShowPostModal(false);
  };

  const handleEventSubmit = async (eventData) => {
    console.log("Submitting eventData:", eventData);
    try {
      const res = await fetch("/frontend-api/events/create", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(eventData),
      });
      if (res.ok) {
        toast.success("Event created!");
        refreshEvents();
      } else {
        toast.error("Failed to create event.");
      }
    } catch (err) {
      toast.error("Network error.");
    }
    setShowEventModal(false);
  };

  // RSVP handler for modal
  const handleRSVP = async (status) => {
    console.log(selectedEvent);
    if (!selectedEvent) return;
    try {
      const res = await fetch("/frontend-api/groups/events/response", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          event_uuid: selectedEvent.uuid,
          response: status,
        }),
      });
      let data = {};
      try {
        data = await res.json();
      } catch (e) {
        // If response is not JSON, data stays empty
      }
      if (res.ok) {
        setCurrentRSVP(status);
        toast.success(
          `RSVP updated: ${status === "accepted" ? "Going" : "Not Going"}`
        );
        refreshEvents();
        fetchAttendees(selectedEvent.uuid); // Refresh attendees after RSVP
      } else {
        console.error("Backend error:", data);
        toast.error("Failed to update RSVP.");
      }
    } catch (err) {
      console.error("Network error:", err);
      toast.error("Network error.");
    }
  };

  return (
    <div>
      <GroupHeader
        group={group}
        isMember={isMember}
        onShowPostModal={() => setShowPostModal(true)}
        onShowEventModal={() => setShowEventModal(true)}
        onRequestJoin={onRequestJoin}
      />
      <div className="group-detail">
        {onBack && (
          <button onClick={onBack} className="group-detail-back-btn">
            ← Back to Groups
          </button>
        )}

      <section className="group-detail-members-section">
          <div
              className="group-members-label group-members-label--interactive"
              onClick={handleOpenMembersModal}
              tabIndex={0}
              role="button"
              aria-label="Show members"
              onKeyDown={e => {
                  if (e.key === "Enter" || e.key === " ") handleOpenMembersModal();
              }}
              >
              <strong>Members:</strong> {group.members_count}
          </div>
          <div className="group-members-actions">
          {isMember ? (
              <button onClick={handleOpenMembersModal}>Invite User</button>
          ) : isRequested ? (
              <button disabled>Pending</button>
          ) : isInvited ? (
              <>
                  <button
                      onClick={() => handleAcceptInvite(currentUserUUID)}
                      style={{ marginRight: '0.5rem' }}
                      disabled={!!loadingActions[currentUserUUID]}
                  >
                      Accept
                  </button>
                  <button
                      onClick={() => handleDeclineInvite(currentUserUUID)}
                      style={{ color: 'red' }}
                      disabled={!!loadingActions[currentUserUUID]}
                  >
                      Decline
                  </button>
              </>
          ) : (
              <button onClick={onRequestJoin} disabled={isJoining}>
                  {isJoining ? 'Pending...' : 'Request to Join'}
              </button>
          )}
          </div>
      </section>

        {/* Events Section */}
        <section className="event-section">
          <strong>Upcoming Events:</strong>
          {loadingEvents ? (
            <div>Loading events...</div>
          ) : groupEvents.length === 0 ? (
            <div>No events yet.</div>
          ) : (
            <ul className="group-detail-events-list">
              {groupEvents.map((event) => (
                <li key={event.uuid} className="group-detail-event-item">
                  <EventCard
                    event={event}
                    isMember={isMember}
                    onClick={() => {
                      setSelectedEvent(event);
                      setCurrentRSVP(event.status);
                      fetchAttendees(event.uuid); // Fetch attendees when opening modal
                    }}
                  />
                </li>
              ))}
            </ul>
          )}
        </section>
        <section className="post-section">
          <strong>Posts:</strong>
          <PostList posts={posts} />
        </section>

        {/* Modals */}
        {showPostModal && (
          <Modal onClose={() => setShowPostModal(false)}>
            <CreatePostForm
              groupID={group.id}
              onSubmit={handlePostSubmit}
              onClose={() => setShowPostModal(false)}
            />
          </Modal>
        )}

        {showEventModal && (
          <Modal onClose={() => setShowEventModal(false)}>
            <CreateEventForm
              groupUUID={group.uuid}
              onSubmit={handleEventSubmit}
              onClose={() => setShowEventModal(false)}
            />
          </Modal>
        )}

        {/* THIS IS YOUR MEMBER DETAIL MODAL! */}
        {showMembersModal && (
            <Modal title="Group Members" onClose={handleCloseMembersModal}>
                <section>
                    <h3>Members</h3>
                    <ul>
                    {acceptedMembers.length > 0 ? (
                            acceptedMembers.map(m => (
                            <li key={m.follower_uuid || m.uuid}>{m.follower_name || m.name}</li>
                            ))
                        ) : (
                            <li>No members</li>
                        )}
                    </ul>
                </section>
                {isMember && (
                <section>
                  <h3>Pending Members</h3>
                  <ul>
                    {pendingRequests.length > 0 ? (
                      pendingRequests.map(r => (
                        <li key={r.follower_uuid || r.uuid || r.user_uuid}>
                          {r.follower_name || r.name}
                          {r.status === "invited" ? (
                            <span style={{ marginLeft: '1rem', color: '#888' }}>Invited</span>
                          ) : (
                            isCreator ? (
                              <>
                                <button
                                  style={{ marginLeft: '1rem' }}
                                  onClick={() => onApproveRequest(r.follower_uuid || r.uuid || r.user_uuid)}
                                  disabled={!!loadingActions[r.follower_uuid || r.uuid || r.user_uuid]}
                                >
                                  {loadingActions[r.follower_uuid || r.uuid || r.user_uuid] ? 'Approving...' : 'Approve'}
                                </button>
                                <button
                                  style={{ marginLeft: '0.5rem', color: 'red' }}
                                  onClick={() => onDenyRequest(r.follower_uuid || r.uuid || r.user_uuid)}
                                  disabled={!!loadingActions[r.follower_uuid || r.uuid || r.user_uuid]}
                                >
                                  {loadingActions[r.follower_uuid || r.uuid || r.user_uuid] ? 'Denying...' : 'Deny'}
                                </button>
                              </>
                            ) : (
                              <span style={{ marginLeft: '1rem', color: '#888' }}>Pending...</span>
                            )
                          )}
                        </li>
                      ))
                    ) : (
                      <li>No pending members</li>
                    )}
                  </ul>
                </section>
                )}
                {isMember && (
                    <section className="invite-non-members">
                        <h3>Non-Members</h3>
                        {loadingUsers ? (
                            <div>Loading users...</div>
                        ) : (
                            <ul>
                            {filteredNonMembers.length > 0 ? (
                                [...filteredNonMembers]
                                .sort((a, b) => a.nick_name.localeCompare(b.nick_name))
                                .map(u => (
                                    <li key={u.uuid}>
                                        {u.nick_name}
                                        <button className="invite-button"
                                            style={{ marginLeft: '1rem' }}
                                            onClick={() => onInviteUser(u.uuid)}
                                            disabled={!!loadingActions[u.uuid]}
                                        >
                                            {loadingActions[u.uuid] ? 'Inviting...' : 'Invite'}
                                        </button>
                                    </li>
                                ))
                            ) : (
                                <li>Everyone is a member!</li>
                            )}
                            </ul>
                        )}
                    </section>
                )}
            </Modal>
        )} 

        {/* Event detail modal with enhanced RSVP logic */}
        {selectedEvent && (
          <Modal
            onClose={() => {
              setSelectedEvent(null);
              setCurrentRSVP(null);
              setEventAttendees([]);
            }}
            title={selectedEvent.title}
          >
            <div className="event-modal-content">
              <div className="event-modal-info">
                <p>
                  <strong>Date/Time:</strong>{" "}
                  {formatDate(selectedEvent.start_time)}
                </p>
                {selectedEvent.location && (
                  <p>
                    <strong>Location:</strong> {selectedEvent.location}
                  </p>
                )}
                <p>
                  <strong>Description:</strong> {selectedEvent.description}
                </p>
              </div>

              <AttendeesDropdown
                eventUUID={selectedEvent.uuid}
                attendCount={selectedEvent.attend_count}
                attendees={eventAttendees}
              />

              {isMember && (
                <div className="rsvp-section">
                  {isUserGoing ? (
                    <>
                      <p className="rsvp-status-message rsvp-status-attending">
                        ✅ You're attending this event
                      </p>
                      <button
                        onClick={() => handleRSVP("declined")}
                        className="rsvp-btn rsvp-btn-not-going"
                      >
                        Not Going
                      </button>
                    </>
                  ) : (
                    <>
                      <p className="rsvp-status-message rsvp-status-not-attending">
                        ❌ You're not attending this event
                      </p>
                      <button
                        onClick={() => handleRSVP("accepted")}
                        className="rsvp-btn rsvp-btn-going"
                      >
                        Going
                      </button>
                    </>
                  )}
                </div>
              )}

            </div>
          </Modal>
        )}
      </div>
    </div>
  );
}
