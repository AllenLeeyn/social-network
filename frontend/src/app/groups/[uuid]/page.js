'use client';

import { useParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import GroupDetail from '../../../components/groups/GroupDetail';
import SidebarSection from '../../../components/SidebarSection';
import UsersList from '../../../components/UsersList';
import GroupMembersList from '../../../components/groups/GroupMemberList';
import '../groups.css'; // or './uuid.css' if you want special styling

export default function GroupDetailPage() {
    const { uuid } = useParams();
    const [group, setGroup] = useState(null);
    const [members, setMembers] = useState([]);
    const [loadingGroup, setLoadingGroup] = useState(true);
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
        fetch(`/frontend-api/groups/members/${uuid}`)
            .then(res => res.json())
            .then(data => {
                console.log("Members data:", data.data)
                setMembers(data.data || []);
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
                        : <GroupMembersList members={members} />
                    }
                </SidebarSection>
            </aside>

            {/* Main Content */}
            <section className="main-feed group-section">
                <GroupDetail group={group} />
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
