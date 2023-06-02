import React, { useEffect, useState } from "react";
import { observer } from "mobx-react-lite";
import UserStore from "../store";
import {listCandidates, listCandidatesVoteStatus, vote} from "../services/api";
import { Table, Button, Message } from "@arco-design/web-react";

interface Candidate {
  Id: number;
  Username: string;
  Voted: boolean;
}

const Vote: React.FC = observer(() => {
  const [candidates, setCandidates] = React.useState<Candidate[]>([]);
  const [voted, setVoted] = useState<boolean>(false);

  useEffect(() => {
    fetchCandidates();
  }, []);

  const fetchCandidates = async () => {
    const candidates = await listCandidatesVoteStatus();
    const list: any = candidates || [];
    console.log("候选人列表", list);

    // 找到一个投过票的，说明已投过，只能投一次
    const _hasVoted = !!list?.find((c:any) => c.Voted);
    if (_hasVoted) {
      Message.info("您已投过票，无法再投");
    }
    setVoted(_hasVoted);
    setCandidates(list);
  };

  const handleVote =  (username: string) => {
    console.log('aaaa vote', username)
   vote(username).then(()=>{
     Message.info("投票成功");
     fetchCandidates();
   }, (e:any) =>{
     console.error('投票失败');
     Message.info("投票失败");
   })

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
      render: (_:any, record:Candidate) => {

        const handleClick = (e:any) => {
          console.log('click vote', e,record);
          handleVote(record.Username)
        }
        return <div onClick={handleClick}>
          <Button
              disabled={voted}
              type="primary"
          >
            投票
          </Button>
        </div>
      },
    },
  ];

  return (
    <div>
      <Table columns={columns} data={candidates} rowKey="Id" />
    </div>
  );
});

export default Vote;
