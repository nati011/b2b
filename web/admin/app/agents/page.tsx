import TableLayout from "../components/table-layout";
import { column, data } from "./column"

export default function Agents() {
  return (
    <TableLayout
       column={column}
       data={data}
        heading={"Agents"}
         pages={[
          {
           "title":"Agents",
            "href":"/agents"
           },
         ]}
         buttonURL="/agents/form"
      />
  );
}
