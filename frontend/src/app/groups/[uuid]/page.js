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

    const [allUsers, setAllUsers] = useState([]);
    const [loadingUsers, setLoadingUsers] = useState(true);



    // Helper to refresh members and requests
    function refreshMembersAndRequests() {
        setLoadingMembers(true);
        Promise.all([
            fetch(`/frontend-api/groups/members/${uuid}`).then(res => res.json()),
            fetch(`/frontend-api/group/member/requests/${uuid}`).then(res => res.json())
        ]).then(([membersData, requestsData]) => {
            setMembers(membersData.data || []);
            setRequests(requestsData.data || []);
            setLoadingMembers(false);
        });
    }

    useEffect(() => {
        setLoadingGroup(true);
        fetch(`/frontend-api/groups/${uuid}`)
        .then(res => res.json())
        .then(data => {
            setGroup(data.data); 
            setLoadingGroup(false);
        });
    }, [uuid]);

    useEffect(() => {
        refreshMembersAndRequests();
    }, [uuid]);


    // fetching all users for group detail comp
    useEffect(() => {
        setLoadingUsers(true);
        fetch('/frontend-api/users')
            .then(res => res.json())
            .then(data => {
                setAllUsers(data.data || []);
                setLoadingUsers(false);
            });
    }, []);

    // Add this function definition:
    function handleRequestJoin() {
        fetch(`/frontend-api/groups/join`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ group_uuid: uuid }),
        })
        .then(res => res.json())
        .then(data => {
            if (data.data) {
                setGroup(data.data);
            } else {
                setGroup(prev => ({ ...prev, status: 'requested' }));
            }
            toast.success('Request Sent!');
            refreshMembersAndRequests();
        })
        .catch(() => {
            toast.error('Failed to send join request.');
        });
    }

    function handleApproveRequest(follower_uuid) {
        fetch(`/frontend-api/group/member/response`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                group_uuid: uuid,
                follower_uuid,
                status: 'accepted'
            }),
        })
        .then(res => res.json())
        .then(data => {
            toast.success('Request approved!');
            refreshMembersAndRequests();
        })
        .catch(() => {
            toast.error('Failed to approve request.');
        });
    }

    function handleDenyRequest(follower_uuid) {
        fetch(`/frontend-api/group/member/response`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                group_uuid: uuid,
                follower_uuid,
                status: 'declined'
            }),
        })
        .then(res => res.json())
        .then(data => {
            toast.info('Request denied.');
            refreshMembersAndRequests();
        })
        .catch(() => {
            toast.error('Failed to deny request.');
        });
    }


    function handleInviteUser(follower_uuid) {
        fetch('/frontend-api/group/invite', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                follower_uuid,
                group_uuid: uuid,
            }),
        })
        .then(res => res.json())
        .then(data => {
            if (data.success) {
                toast.success('User invited!');
                refreshMembersAndRequests();
            } else {
                toast.error(data.message || 'Failed to invite user.');
            }
        })
        .catch(() => {
            toast.error('Failed to invite user.');
        });
    }

    // handlers for accepting and declining as user who has been invited
    function handleAcceptInvite(follower_uuid) {
        fetch(`/frontend-api/group/member/response`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                group_uuid: uuid,
                follower_uuid,
                status: 'accepted'
            }),
        })
        .then(res => res.json())
        .then(data => {
            toast.success('Invite accepted!');
            refreshMembersAndRequests();
        })
        .catch(() => {
            toast.error('Failed to accept invite.');
        });
    }
    function handleDeclineInvite(follower_uuid) {
        fetch(`/frontend-api/group/member/response`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
                group_uuid: uuid,
                follower_uuid,
                status: 'declined'
            }),
        })
        .then(res => res.json())
        .then(data => {
            toast.info('Invite declined.');
            refreshMembersAndRequests();
        })
        .catch(() => {
            toast.error('Failed to decline invite.');
        });
    }



    if (loadingGroup) return <div>Loading...</div>;
    if (!group) return <div>Group not found.</div>;


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
                        : <GroupMembersList 
                            members={members.filter(m => m.status === 'accepted')} 
                            requests={requests} 
                            groupUuid={uuid}
                        />
                    }
                </SidebarSection>
            </aside>

            {/* Main Content */}
            <section className="main-feed group-section">
                <GroupDetail 
                    group={group}
                    members={members}
                    requests={requests} 
                    allUsers={allUsers}
                    loadingUsers={loadingUsers}
                    onRequestJoin={handleRequestJoin} 
                    onApproveRequest={handleApproveRequest}
                    onDenyRequest={handleDenyRequest}
                    onInviteUser={handleInviteUser}
                    handleAcceptInvite={handleAcceptInvite}
                    handleDeclineInvite={handleDeclineInvite}
                />
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
