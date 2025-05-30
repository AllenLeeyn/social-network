'use client';

import { useParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import GroupDetail from '../../../components/groups/GroupDetail';
import SidebarSection from '../../../components/SidebarSection';
import UsersList from '../../../components/UsersList';
import '../groups.css'; // or './uuid.css' if you want special styling

export default function GroupDetailPage() {
    const { uuid } = useParams();
    const [group, setGroup] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetch(`/frontend-api/groups/${uuid}`)
        .then(res => res.json())
        .then(data => {
            setGroup(data.data); 
            setLoading(false);
        });
    }, [uuid]);

    if (loading) return <div>Loading...</div>;
    if (!group) return <div>Group not found.</div>;

    return (
        <div className="groups-page-layout">
            {/* Left Sidebar */}
            <aside className="sidebar left-sidebar">
                {/* You can put group navigation, filters, or actions here */}
                <SidebarSection title="Group Actions">
                    {/* Example: */}
                    {/* <button>Invite to Group</button> */}
                    {/* <button>Request to Join</button> */}
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
