'use client';

import { useParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { fetchGroupPosts } from "../../../lib/apiPosts";
import GroupDetail from '../../../components/groups/GroupDetail';
import SidebarSection from '../../../components/SidebarSection';
import UsersList from '../../../components/UsersList';
import GroupMembersList from '../../../components/groups/GroupMemberList';
import '../groups.css'; 

export default function GroupDetailPage() {
    const { uuid } = useParams();

    const [group, setGroup] = useState(null);
    const [loadingGroup, setLoadingGroup] = useState(true);

    const [members, setMembers] = useState([]);
    const [requests, setRequests] = useState([]);
    const [loadingMembers, setLoadingMembers] = useState(true);

    const [allUsers, setAllUsers] = useState([]);
    const [loadingUsers, setLoadingUsers] = useState(true);

    const [isJoining, setIsJoining] = useState(false);
    const [isInviting, setIsInviting] = useState(false);
    const [loadingActions, setLoadingActions] = useState({});

    function setActionLoading(uuid, value) {
        setLoadingActions(prev => ({ ...prev, [uuid]: value }));
    }

    function refreshGroup() {
        setLoadingGroup(true);
        fetch(`/frontend-api/group/${uuid}`)
            .then(res => res.json())
            .then(data => {
                setGroup(data.data);
                setLoadingGroup(false);
            });
    }
    const [grpPosts, setGrpPosts] = useState([]);
    useEffect(() => {

        async function loadGroupPost() {
            try {
            const groupPostsData = await fetchGroupPosts(uuid);
            setGrpPosts(groupPostsData.data);

            } catch (err) {
                toast.error(err.message);
            }
        }
        loadGroupPost();
    }, [uuid]);

    async function refreshMembersAndRequests() {
        setLoadingMembers(true);
        try {
            const [membersRes, requestsRes] = await Promise.all([
                fetch(`/frontend-api/groups/members/${uuid}`),
                fetch(`/frontend-api/group/member/requests/${uuid}`)
            ]);
            const membersData = await membersRes.json();
            const requestsData = await requestsRes.json();
            setMembers(membersData.data || []);
            setRequests(requestsData.data || []);
        } finally {
            setLoadingMembers(false);
        }
    }


    useEffect(() => {
        refreshGroup();
    }, [uuid]);

    useEffect(() => {
        refreshMembersAndRequests();
    }, [uuid]);

    // useEffect(() => {
    //     setLoadingGroup(true);
    //     fetch(`/frontend-api/groups/${uuid}`)
    //     .then(res => res.json())
    //     .then(data => {
    //         setAllUsers(data.data || []);
    //         setLoadingUsers(false);
    //     });
    // }, [uuid]);


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

function handleRequestJoin() {
        if (isJoining) return;
        setIsJoining(true);
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
            refreshGroup();
        })
        .catch(() => toast.error('Failed to send join request.'))
        .finally(() => setIsJoining(false));
    }

    function handleApproveRequest(follower_uuid) {
        if (loadingActions[follower_uuid]) return;
        setActionLoading(follower_uuid, true);
        fetch(`/frontend-api/group/member/response`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ group_uuid: uuid, follower_uuid, status: 'accepted' }),
        })
        .then(res => res.json())
        .then(() => {
            toast.success('Request approved!');
            refreshMembersAndRequests();
        })
        .catch(() => toast.error('Failed to approve request.'))
        .finally(() => setActionLoading(follower_uuid, false));
    }

    function handleDenyRequest(follower_uuid) {
        if (loadingActions[follower_uuid]) return;
        setActionLoading(follower_uuid, true);
        fetch(`/frontend-api/group/member/response`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ group_uuid: uuid, follower_uuid, status: 'declined' }),
        })
        .then(res => res.json())
        .then(() => {
            toast.info('Request denied.');
            refreshMembersAndRequests();
        })
        .catch(() => toast.error('Failed to deny request.'))
        .finally(() => setActionLoading(follower_uuid, false));
    }

    function handleInviteUser(follower_uuid) {
        if (loadingActions[follower_uuid]) return;
        setActionLoading(follower_uuid, true);
        fetch('/frontend-api/group/invite', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ follower_uuid, group_uuid: uuid }),
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
        .catch(() => toast.error('Failed to invite user.'))
        .finally(() => {
            setActionLoading(follower_uuid, false);
        });
    }

    // handlers for accepting and declining as user who has been invited
    function handleAcceptInvite(follower_uuid) {
        if (loadingActions[follower_uuid]) return;
        setActionLoading(follower_uuid, true);
        fetch(`/frontend-api/group/member/response`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ group_uuid: uuid, follower_uuid, status: 'accepted' }),
        })
        .then(res => res.json())
        .then(() => {
            toast.success('Invite accepted!');
            refreshMembersAndRequests();
            refreshGroup();
        })
        .catch(() => toast.error('Failed to accept invite.'))
        .finally(() => setActionLoading(follower_uuid, false));
    }

    function handleDeclineInvite(follower_uuid) {
        if (loadingActions[follower_uuid]) return;
        setActionLoading(follower_uuid, true);
        fetch(`/frontend-api/group/member/response`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ group_uuid: uuid, follower_uuid, status: 'declined' }),
        })
        .then(res => res.json())
        .then(() => {
            toast.info('Invite declined.');
            refreshMembersAndRequests();
        })
        .catch(() => toast.error('Failed to decline invite.'))
        .finally(() => setActionLoading(follower_uuid, false));
    }

    if (loadingGroup) return <div>Loading...</div>;
    if (!group) return <div>Group not found.</div>;


    return (
        <div className="groups-page-layout">
            <aside className="sidebar left-sidebar">
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
                    loadingActions={loadingActions}
                    isJoining={isJoining}
                    isInviting={isInviting}
                    posts={grpPosts}
                />
            </section>

            {/* Right Sidebar */}
            <aside className="sidebar right-sidebar">
                <SidebarSection title="Chat list">
                    <UsersList />
                </SidebarSection>
            </aside>
        </div>
    );
}
