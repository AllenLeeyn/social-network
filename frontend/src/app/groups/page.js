// app/groups/page.js

'use client';

import React, { useState } from 'react';
import SidebarSection from '../../components/SidebarSection';
import UsersList from '../../components/UsersList';
import GroupFilterList from '../../components/groups/GroupFilterList';
import GroupList from '../../components/groups/GroupList';
import GroupDetail from '../../components/groups/GroupDetail';
import GroupInvitationList from '../../components/groups/GroupInvitationList';

const groupFilters = [
  { key: 'my_groups', label: 'My Groups' },
  { key: 'discover', label: 'Discover' },
  { key: 'invitations', label: 'Invitations' }
];

export default function GroupsPage() {
  const [selectedFilter, setSelectedFilter] = useState('my_groups');
  const [selectedGroup, setSelectedGroup] = useState(null);

  return (
    <main>
      <div className='groups-page-layout'>
        {/* Left Sidebar */}
        <aside className='sidebar left-sidebar'>
          <SidebarSection title='Groups'>
            <GroupFilterList
              filters={groupFilters}
              selectedFilter={selectedFilter}
              onSelect={setSelectedFilter}
            />
          </SidebarSection>
        </aside>

        {/* Main Content */}
        <section className='main-feed group-section'>
          {selectedFilter === 'my_groups' && (
            <GroupList onSelectGroup={setSelectedGroup} />
          )}
          {selectedFilter === 'discover' && (
            <GroupList discover onSelectGroup={setSelectedGroup} />
          )}
          {selectedFilter === 'invitations' && (
            <GroupInvitationList />
          )}
          {selectedGroup && (
            <GroupDetail group={selectedGroup} />
          )}
        </section>

        {/* Right Sidebar */}
        <aside className="sidebar right-sidebar">
          <SidebarSection title="All Users">
            <UsersList />
          </SidebarSection>
        </aside>
      </div>
    </main>
  );
}
