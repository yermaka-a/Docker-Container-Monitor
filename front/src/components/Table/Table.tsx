import 'bootstrap/dist/css/bootstrap.min.css';
import { useEffect, useState } from 'react';
import {Container} from './types.ts'


const Table: React.FC = () => {


    const [containers, setContainers] = useState<Container[] | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
  
    async function fetchData<T>(): Promise<T[]>{
        const url = "http://localhost:3001/containers"
        const response = await fetch(url)
        if (!response.ok) {
          throw new Error("Network response wasn't ok")
        } 
        return response.json() 
  }
    const getContainers = async () => {
      try {
          const containers = await fetchData<Container>()
          console.log(containers)
          if (containers == null){
            setContainers(containers)
            return
          }
          containers.sort((prev, next)=>prev.id - next.id)
          let i = 0
          containers.map((el)=> {
            i+=1
            return el.id = i
          })
          setContainers(containers)
        } 
       catch (error) {
        if (error instanceof Error){
        setError(error.message)
        setContainers(null)
        }
      } finally {
        setLoading(false)
      }
  }

    useEffect(() => {
      getContainers()
      const intervalId = setInterval(getContainers, 3000) 
  
      return () => clearInterval(intervalId)
    }, []); 
    if (containers !== null)  return (
      <div className="container">
        <table className="table table-striped">
          <thead className="thead-dark">
            <tr>
            <th scope="col">ID</th>
              <th scope="col">IP</th>
              <th scope="col">Время пинга в ms</th>
              <th scope="col">Дата последней успешной попытки</th>
            </tr>
          </thead>
          <tbody>
            {containers.map((item, index) => (
              <tr key={index}>
                <td>{item.id}</td>
                <td>{item.ip}</td>
                <td>{item.timeMs}</td>
                <td>{item.pingDate}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    )
    else if (loading || null) return (<div className="container">
      <table className="table table-striped">
        <thead className="thead-dark">
          <tr>
          <th scope="col">ID</th>
            <th scope="col">IP</th>
            <th scope="col">Время пинга в ms</th>
            <th scope="col">Дата последней успешной попытки</th>
          </tr>
        </thead>
        <tbody>
            <tr key={0}>
              <td colSpan={4}>{"В данный момент нет опрошенных контейнеров..."}</td>
            </tr>
        </tbody>
      </table>
    </div>)
    if (error) return (
    <div className="container">
    <table className="table table-striped">
        <thead className="thead-dark">
          <tr>
          <th scope="col">Возникла ошибка</th>
          </tr>
        </thead>
        <tbody>
            <tr key={0}>
              <td colSpan={4}>{error}</td>
            </tr>
        </tbody>
      </table>
      </div>
    )
 
};

export default Table;