import React, { useEffect, useState } from "react";
import { observer } from "mobx-react-lite";
import UserStore from "../store";
import { listCandidates } from "../services/api";
import { Table, Button, Message } from "@arco-design/web-react";

interface Candidate {
  Id: number;
  Username: string;
  voted: boolean;
}

const Vote: React.FC = observer(() => {
  const [candidates, setCandidates] = React.useState<Candidate[]>([]);
  const [voted, setVoted] = useState<boolean>(false);

  useEffect(() => {
    fetchCandidates();
  }, []);

  const fetchCandidates = async () => {
    const candidates = await listCandidates();
    const list: Candidate[] = candidates?.data || [];
    console.log("候选人列表", list);

    // 找到一个投过票的，说明已投过，只能投一次
    const _hasVoted = !!list?.find((c) => c.voted);
    if (_hasVoted) {
      Message.info("您已投过票，无法再投");
    }
    setVoted(_hasVoted);
    setCandidates(list);
  };

  const handleVote = (username: string) => {
    // 投票
    console.log("投票", username);
  };

  const columns = [
    {
      title: "ID",
      dataIndex: "Id",
      key: "Id",
    },
    {
      title: "用户名",
      dataIndex: "Username",
      key: "Username",
    },
    {
      title: "投票",
      key: "action",
      render: (record: Candidate) => (
        <>
          {UserStore.username === "admin" && (
            <Button
              disabled={voted}
              onClick={() => handleVote(record.Username)}
            >
              投票
            </Button>
          )}
        </>
      ),
    },
  ];

  return (
    <div>
      <Table columns={columns} data={candidates} rowKey="id" />
    </div>
  );
});

export default Vote;
