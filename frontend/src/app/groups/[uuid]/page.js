'use client';

import { useParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import GroupDetail from '../../../components/groups/GroupDetail';
import SidebarSection from '../../../components/SidebarSection';
import UsersList from '../../../components/UsersList';
import GroupMembersList from '../../../components/groups/GroupMemberList';
import '../groups.css'; // or './uuid.css' if you want special styling

export default function GroupDetailPage() {
    const { uuid } = useParams();
    const [group, setGroup] = useState(null);
    const [loadingGroup, setLoadingGroup] = useState(true);

    const [members, setMembers] = useState([]);
    const [requests, setRequests] = useState([]);
    const [loadingMembers, setLoadingMembers] = useState(true);

    useEffect(() => {
        setLoadingGroup(true);
        fetch(`/frontend-api/groups/${uuid}`)
        .then(res => res.json())
        .then(data => {
            setGroup(data.data); 
            setLoadingGroup(false);
        });
    }, [uuid]);

    // Fetch group members (using your backend handler)
    useEffect(() => {
        setLoadingMembers(true);
        Promise.all([
            fetch(`/frontend-api/groups/members/${uuid}`).then(res => res.json()),
            fetch(`/frontend-api/group/member/requests/${uuid}`).then(res => res.json())
        ]).then(([membersData, requestsData]) => {
            console.log(requestsData)
            setMembers(membersData.data || []);
            setRequests(requestsData.data || []);
            setLoadingMembers(false);
        });
    }, [uuid]);

    if (loadingGroup) return <div>Loading...</div>;
    if (!group) return <div>Group not found.</div>;

    function handleInviteUser() {
    // You can open a modal, show a toast, or implement your invite logic here
    alert("Invite User clicked!");
    // Or setShowInviteModal(true) if you want to open a modal
    }

    async function handleRequestJoin(group) {
        try {
            const response = await fetch('/frontend-api/groups/join', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ group_uuid: group.uuid }),
            });

            // Check for HTTP errors
            if (!response.ok) {
                // Try to parse error message from server, fallback to status text
                let errorMsg = `Request failed (${response.status})`;
                try {
                    const errorData = await response.json();
                    errorMsg = errorData.error || errorMsg;
                } catch {
                    // ignore JSON parse errors
                }
                toast.error(errorMsg);
            }

            const data = await response.json();
            if (data.success) {
                toast.success('Request sent!');
                // Optionally update group status in state here
            } else {
                toast.error(data.error || 'Request failed.');
            }
        } catch (error) {
            // Handles network errors, timeouts, etc.
            toast.error('Network error. Please try again.');
        }
    }


    return (
        <div className="groups-page-layout">
            {/* Left Sidebar */}
            <aside className="sidebar left-sidebar">
                <SidebarSection title="Group Actions">
                    <button onClick={handleInviteUser}>Invite User</button>
                    {/* <button>Request to Join</button> */}
                </SidebarSection>
                {/* Members List in Sidebar */}
                <SidebarSection title="Members">
                    {loadingMembers
                        ? <div>Loading members...</div>
                        : <GroupMembersList members={members} requests={requests} />
                    }
                </SidebarSection>
            </aside>

            {/* Main Content */}
            <section className="main-feed group-section">
                <GroupDetail group={group} onRequestJoin={handleRequestJoin} />
            </section>

            {/* Right Sidebar */}
            <aside className="sidebar right-sidebar">
                <SidebarSection title="All Users">
                    <UsersList />
                </SidebarSection>
            </aside>
        </div>
    );
}
